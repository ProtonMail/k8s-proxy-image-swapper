// SPDX-License-Identifier: Apache-2.0
package mutate

import "testing"

type TestImageCase struct {
	image    string
	registry string
	expected string
}

func TestConcatenation(t *testing.T) {
    arr := []string{"toto", "tata"}

    result := concatenateStringArray(arr, ":")
    if result != "toto:tata" {
        t.Errorf("concatenation is not working properly")
    }
}

func TestImage(t *testing.T) {
	testTable := []TestImageCase{
		{
			image:    "busybox:toto",
			registry: "example.com",
			expected: "example.com/docker.io/library/busybox:toto",
		},
		{
			image:    "busybox",
			registry: "example.com",
			expected: "example.com/docker.io/library/busybox:latest",
		},
		{
			image:    "toto/tata:titi",
			registry: "example.com",
			expected: "example.com/docker.io/toto/tata:titi",
		},
		{
			image:    "toto/tata",
			registry: "example.com",
			expected: "example.com/docker.io/toto/tata:latest",
		},
		{
			image:    "example.org/toto/tata:titi",
			registry: "example.com",
			expected: "example.com/example.org/toto/tata:titi",
		},
		{
			image:    "example.org/toto/tata",
			registry: "example.com",
			expected: "example.com/example.org/toto/tata:latest",
		},
		{
			image:    "docker.io/toto/tata:titi",
			registry: "example.com",
			expected: "example.com/docker.io/toto/tata:titi",
		},
		{
			image:    "docker.io/toto/tata",
			registry: "example.com",
			expected: "example.com/docker.io/toto/tata:latest",
		},
		{
			image:    "docker.io/busybox:titi",
			registry: "example.com",
			expected: "example.com/docker.io/library/busybox:titi",
		},
		{
			image:    "docker.io/busybox",
			registry: "example.com",
			expected: "example.com/docker.io/library/busybox:latest",
		},
		{
			image:    "gcr.io/busybox:titi",
			registry: "example.com",
			expected: "example.com/gcr.io/busybox:titi",
		},
		{
			image:    "gcr.io/toto/tata:titi",
			registry: "example.com",
			expected: "example.com/gcr.io/toto/tata:titi",
		},
		{
			image:    "gcr.io/toto/tata@sha256:XXXXX",
			registry: "example.com",
			expected: "example.com/gcr.io/toto/tata@sha256:XXXXX",
		},
		{
			image:    "gcr.io/toto/tata:titi@sha256:XXXXX",
			registry: "example.com",
			expected: "example.com/gcr.io/toto/tata:titi@sha256:XXXXX",
		},
		{
			image:    "gcr.io/toto/tata/titi:tag",
			registry: "example.com",
			expected: "example.com/gcr.io/toto/tata/titi:tag",
		},
		{
			image:    "registry.org/toto/tata/titi:tag",
			registry: "example.com",
			expected: "example.com/registry.org/toto/tata/titi:tag",
		},
		{
			image:    "toto/tata/titi:tag",
			registry: "example.com",
			expected: "example.com/docker.io/toto/tata/titi:tag",
		},
		{
			image: "quay.io/toto/tata",
			registry: "example.com",
			expected: "example.com/quay.io/toto/tata:latest",
		},
	}

	for _, test := range testTable {
		res := GetPatchedImageUrl(test.image, test.registry)


		if res != test.expected {

			t.Errorf("Error test image : %v, registry : %v, expected : %v, got %v\n",
				test.image,
				test.registry,
				test.expected,
				res)
		}
	}
}
