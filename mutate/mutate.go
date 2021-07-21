// SPDX-License-Identifier: AGPL-3.0-only
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

func GetPatchedImageUrl(img, registry string) string {
	if strings.HasPrefix(img, "registry:") ||
	   strings.Contains(img, "/registry:") {
		return img
	}

	imgArr := strings.Split(img, "/")
	// Not prefixed with a site
	if len(imgArr) == 1 {
		// Case busybox or busybox:tag
		return fmt.Sprintf("%s/library/%s", registry, img)
	}

	imgUrl := imgArr[0]
	imgName := strings.Join(imgArr[1:], "/")
	// Case docker.io/busybox
	if len(imgArr) == 2 && imgUrl == "docker.io" {
		return fmt.Sprintf("%s/library/%s", registry, imgName)
	}

	// Case toto/tata (and ! gcr.io/toto)
	if len(imgArr) == 2 && !strings.Contains(imgUrl, ".") {
		return fmt.Sprintf("%s/%s", registry, strings.Join(imgArr, "/"))
	}

	// Case docker.io/toto/tata
	if imgUrl == "docker.io" {
		return fmt.Sprintf("%s/%s", registry, imgName)
	}
	return img
}

func getPatchFromContainerList(ctn []corev1.Container, registry string) []map[string]string {
	patchList := []map[string]string{}
	for i := range ctn {
		img := ctn[i].Image

		patchedImg := GetPatchedImageUrl(img, registry)

		// In case there's a tag
		if strings.HasPrefix(patchedImg, "docker.io/library/registry") ||
		   strings.HasPrefix(patchedImg, "registry") {
			// We don't patch the registry to avoid the bootstrap problem
			continue
		}

		// No need to patch if it's the same
		if img == patchedImg {
			continue
		}

		patch := map[string]string{
			"op":    "replace",
			"path":  fmt.Sprintf("/spec/containers/%d/image", i),
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
		return nil, fmt.Errorf("Unmarshaling request error %s", err)
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
		log.Println("FATAL Error ", err)
		return nil, fmt.Errorf("Unmarshal pod json error %v", err)
	}

	resp.Allowed = true
	resp.UID = ar.UID
	pT := v1beta1.PatchTypeJSONPatch
	resp.PatchType = &pT

	resp.AuditAnnotations = map[string]string{
		"k8s-proxy-image-swapper": "mutated",
	}

	patchList := getPatchFromContainerList(pod.Spec.Containers, registry)
	patchList = append(patchList, getPatchFromContainerList(pod.Spec.InitContainers, registry)...)
	resp.Patch, err = json.Marshal(patchList)

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
