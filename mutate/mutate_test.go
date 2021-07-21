// SPDX-License-Identifier: AGPL-3.0-only
package mutate

import "testing"
import "fmt"

func TestLibraryImage(t *testing.T) {
	if GetPatchedImageUrl("busybox", "example.com") != "example.com/library/busybox" {
		t.Log("Error TestLibraryImage")
		t.Fail()
	}
}

func TestLibraryImageWithTag(t *testing.T) {
	if GetPatchedImageUrl("busybox:latest", "example.com") != "example.com/library/busybox:latest" {
		t.Log("Error TestLibraryImageWithTag")
		t.Fail()
	}
}

func TestNonLibraryImage(t *testing.T) {
	if GetPatchedImageUrl("toto/tata", "example.com") != "example.com/toto/tata" {
		t.Log("Error TestNonLibraryImage")
		t.Fail()
	}
}

func TestNonLibraryImageWithTag(t *testing.T) {
	if GetPatchedImageUrl("toto/tata:latest", "example.com") != "example.com/toto/tata:latest" {
		t.Log("Error TestNonLibraryImageWithTag")
		t.Fail()
	}
}

func TestFullPathImage(t *testing.T) {
	if GetPatchedImageUrl("example.org/toto/tata", "example.com") != "example.org/toto/tata" {
		t.Log("Error TestFullPathImage")
		t.Fail()
	}
}

func TestFullPathImageWithTag(t *testing.T) {
	if GetPatchedImageUrl("example.org/toto/tata:latest", "example.com") != "example.org/toto/tata:latest" {
		t.Log("Error TestFullPathImageWithTag")
		t.Fail()
	}
}

func TestFullPathImageFromDocker(t *testing.T) {
	if GetPatchedImageUrl("docker.io/toto/tata", "example.com") != "example.com/toto/tata" {
		t.Log("Error TestFullPathImageFromDocker")
		t.Fail()
	}
}

func TestFullPathImageFromDockerWithTag(t *testing.T) {
	if GetPatchedImageUrl("docker.io/toto/tata:latest", "example.com") != "example.com/toto/tata:latest" {
		t.Log("Error TestFullPathImageFromDockerWithTag")
		t.Fail()
	}
}

func TestFullPathImageLibraryFromDocker(t *testing.T) {
	if GetPatchedImageUrl("docker.io/busybox", "example.com") != "example.com/library/busybox" {
		t.Log("Error TestFullPathImageLibraryFromDocker")
		t.Fail()
	}
}

func TestFullPathImageLibraryFromDockerWithTag(t *testing.T) {
	if GetPatchedImageUrl("docker.io/busybox:latest", "example.com") != "example.com/library/busybox:latest" {
		t.Log("Error TestFullPathImageLibraryFromDockerWithTag")
		t.Fail()
	}
}

func TestForeignImage(t *testing.T) {
	fmt.Println(GetPatchedImageUrl("gcr.io/busybox", "example.com"))
	if GetPatchedImageUrl("gcr.io/busybox", "example.com") != "gcr.io/busybox" {
		t.Log("Error TestForeignImage")
		t.Fail()
	}
}

func TestForeignFullPathImage(t *testing.T) {
	if GetPatchedImageUrl("gcr.io/toto/tata", "example.com") != "gcr.io/toto/tata" {
		t.Log("Error TestForeignFullPathImage")
		t.Fail()
	}
}
