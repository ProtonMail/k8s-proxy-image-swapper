// SPDX-License-Identifier: Apache-2.0
package mutate

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// registry/namespace/image:tag
type dockerImageUrl struct {
	registry  string
	namespace string
	image     string
	tag       string
}

// takes toto:tata or toto, gives tata or latest
func getImgTag(img string) string {
	imgTagArr := strings.Split(img, ":")
	tag := "latest"

	if len(imgTagArr) != 1 {
		tag = ""
		for _, v := range(imgTagArr[1:]) {
			tag += v + ":"
		}
		// remove last ":"
		tag = tag[:len(tag) - 1]
	}

	return tag
}

// takes toto:tata or toto, gives toto
func getImgName(img string) string {
	imgTagArr := strings.Split(img, ":")
	return imgTagArr[0]
}

func concatenateStringArray(arr []string, separator string) string {
    result := ""

    for _, v := range arr {
        result += v + separator
    }
    // remove last separator
    result = result[: len(result) - len(separator)]
    return result
}

func getDockerImageUrl(img string) dockerImageUrl {
	imgArr := strings.Split(img, "/")
	// Not prefixed with a site
	if len(imgArr) == 1 {
		// Case busybox or busybox:tag
		return dockerImageUrl{
			registry:  "docker.io", // default
			namespace: "library",   // default
			image:     getImgName(img),
			tag:       getImgTag(img),
		}
	}

	imgUrl := imgArr[0]
	// Case docker.io/busybox
	if len(imgArr) == 2 && imgUrl == "docker.io" {
		return dockerImageUrl{
			registry:  imgUrl,
			namespace: "library",
			image:     getImgName(imgArr[1]),
			tag:       getImgTag(imgArr[1]),
		}
	}

	// Case toto/tata (and ! gcr.io/toto)
	if len(imgArr) == 2 && !strings.Contains(imgUrl, ".") {
		return dockerImageUrl{
			registry:  "docker.io",
			namespace: imgArr[0],
			image:     getImgName(imgArr[1]),
			tag:       getImgTag(imgArr[1]),
		}
	}

	if len(imgArr) == 2 && strings.Contains(imgUrl, ".") {
		return dockerImageUrl{
			registry:  imgUrl,
			namespace: "", // ??? TODO does it exist?
			image:     getImgName(imgArr[1]),
			tag:       getImgTag(imgArr[1]),
		}
	}

	// case toto.io/tata/titi[:tag]
	// or case toto.io/tata/titi/toto[:tag]
	if strings.Contains(imgUrl, ".") {
		return dockerImageUrl{
			registry:  imgArr[0],
		        namespace: concatenateStringArray(imgArr[1: len(imgArr) - 1], "/"),
			image:     getImgName(imgArr[len(imgArr) - 1]),
			tag:       getImgTag(imgArr[len(imgArr) - 1]),
		}
	} else {
		// case toto/tata/titi:tag
		return dockerImageUrl{
			registry:  "docker.io",
			namespace: concatenateStringArray(imgArr[: len(imgArr) - 1], "/"),
			image:     getImgName(imgArr[len(imgArr) - 1]),
			tag:       getImgTag(imgArr[len(imgArr) - 1]),
		}
    }
}

func (i dockerImageUrl) String() string {
	if i.namespace == "" {
		return fmt.Sprintf("%s/%s:%s",
			i.registry,
			i.image,
			i.tag)
	}
	return fmt.Sprintf("%s/%s/%s:%s",
		i.registry,
		i.namespace,
		i.image,
		i.tag)
}

func isSameImage(image1, image2 dockerImageUrl) bool {
	return image1.registry == image2.registry &&
		image1.namespace == image2.namespace &&
		image1.image == image2.image
}

func GetPatchedImageUrl(img, registry string) string {
	patchimg := getDockerImageUrl(img)

	for _, image := range Configuration.IgnoreImages {
		if getDockerImageUrl(image).String() == patchimg.String() ||
			isSameImage(getDockerImageUrl(image), patchimg) {
			return patchimg.String()
		}
	}

	if patchimg.registry == "docker.io" {
		patchimg.registry = registry
	}

	return patchimg.String()
}

func getPatchFromContainerList(ctn []corev1.Container, registry, containerType string) []map[string]string {
	patchList := []map[string]string{}
	for i := range ctn {
		img := ctn[i].Image

		patchedImg := GetPatchedImageUrl(img, registry)

		// No need to patch if it's the same
		if img == patchedImg {
			continue
		}

		patch := map[string]string{
			"op":    "replace",
			"path":  fmt.Sprintf("/spec/%s/%d/image", containerType, i),
			"value": patchedImg,
		}
		patchList = append(patchList, patch)
	}

	return patchList
}

func Mutate(body []byte, verbose bool, registry string) ([]byte, error) {
	if verbose {
		log.Printf("recv: %s\n", string(body))
	}

	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("Unmarshaling request error %w", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}

	if ar == nil {
		if verbose {
			log.Printf("resp: %s\n", string(responseBody))
		}

		return responseBody, nil
	}

	if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
		log.Println("Unmarshal pod json error", err)
		return nil, fmt.Errorf("Unmarshal pod json error %w", err)
	}

	resp.Allowed = true
	resp.UID = ar.UID
	pT := v1beta1.PatchTypeJSONPatch
	resp.PatchType = &pT

	resp.AuditAnnotations = map[string]string{
		"k8s-proxy-image-swapper": "mutated",
	}

	patchList := getPatchFromContainerList(pod.Spec.Containers, registry, "containers")
	patchList = append(patchList, getPatchFromContainerList(pod.Spec.InitContainers, registry, "initContainers")...)
	annotationsPatch := map[string]string{
		"op":    "add",
		"path":  "/metadata/labels/k8s-proxy-image-swapper",
		"value": "patched-image",
	}
	if len(patchList) != 0 {
		patchList = append(patchList, annotationsPatch)
	}

	resp.Patch, err = json.Marshal(patchList)
	if err != nil {
		log.Println("Error unmarshalling patchList into AdmissionResponse", err)
		return nil, fmt.Errorf("Error unmarshalling patchList into AdmissionResponse %w", err)
	}

	// We cannot fail
	resp.Result = &metav1.Status{
		Status: "Success",
	}

	admReview.Response = &resp
	responseBody, err = json.Marshal(admReview)
	if err != nil {
		log.Println("FATAL Error ", err)
		return nil, err
	}
	return responseBody, nil
}
