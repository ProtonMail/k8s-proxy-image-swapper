// SPDX-License-Identifier: Apache-2.0
package mutate

import "testing"

type TestImageCase struct {
	image    string
	registry string
	expected string
}

func TestImage(t *testing.T) {
	testTable := []TestImageCase{
		{
			image:    "busybox:toto",
			registry: "example.com",
			expected: "example.com/library/busybox:toto",
		},
		{
			image:    "busybox",
			registry: "example.com",
			expected: "example.com/library/busybox:latest",
		},
		{
			image:    "toto/tata:titi",
			registry: "example.com",
			expected: "example.com/toto/tata:titi",
		},
		{
			image:    "toto/tata",
			registry: "example.com",
			expected: "example.com/toto/tata:latest",
		},
		{
			image:    "example.org/toto/tata:titi",
			registry: "example.com",
			expected: "example.org/toto/tata:titi",
		},
		{
			image:    "example.org/toto/tata",
			registry: "example.com",
			expected: "example.org/toto/tata:latest",
		},
		{
			image:    "docker.io/toto/tata:titi",
			registry: "example.com",
			expected: "example.com/toto/tata:titi",
		},
		{
			image:    "docker.io/toto/tata",
			registry: "example.com",
			expected: "example.com/toto/tata:latest",
		},
		{
			image:    "docker.io/busybox:titi",
			registry: "example.com",
			expected: "example.com/library/busybox:titi",
		},
		{
			image:    "docker.io/busybox",
			registry: "example.com",
			expected: "example.com/library/busybox:latest",
		},
		{
			image:    "gcr.io/busybox:titi",
			registry: "example.com",
			expected: "gcr.io/busybox:titi",
		},
		{
			image:    "gcr.io/toto/tata:titi",
			registry: "example.com",
			expected: "gcr.io/toto/tata:titi",
		},
		{
			image:    "gcr.io/toto/tata:titi@sha256:XXXXX",
			registry: "example.com",
			expected: "gcr.io/toto/tata:titi@sha256:XXXXX",
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
