// SPDX-License-Identifier: Apache-2.0
package mutate

import "testing"

func TestLibraryImage(t *testing.T) {
	if getPatchedImageUrl("busybox:toto", "example.com") != "example.com/library/busybox:toto" {
		t.Log("Error TestLibraryImage")
		t.Fail()
	}
}

func TestLibraryImageWithTag(t *testing.T) {
	if getPatchedImageUrl("busybox", "example.com") != "example.com/library/busybox:latest" {
		t.Log("Error TestLibraryImageWithTag")
		t.Fail()
	}
}

func TestNonLibraryImage(t *testing.T) {
	if getPatchedImageUrl("toto/tata:titi", "example.com") != "example.com/toto/tata:titi" {
		t.Log("Error TestNonLibraryImage")
		t.Fail()
	}
}

func TestNonLibraryImageWithTag(t *testing.T) {
	if getPatchedImageUrl("toto/tata", "example.com") != "example.com/toto/tata:latest" {
		t.Log("Error TestNonLibraryImageWithTag")
		t.Fail()
	}
}

func TestFullPathImage(t *testing.T) {
	if getPatchedImageUrl("example.org/toto/tata:titi", "example.com") != "example.org/toto/tata:titi" {
		t.Log("Error TestFullPathImage")
		t.Fail()
	}
}

func TestFullPathImageWithTag(t *testing.T) {
	if getPatchedImageUrl("example.org/toto/tata", "example.com") != "example.org/toto/tata:latest" {
		t.Log("Error TestFullPathImageWithTag")
		t.Fail()
	}
}

func TestFullPathImageFromDocker(t *testing.T) {
	if getPatchedImageUrl("docker.io/toto/tata:titi", "example.com") != "example.com/toto/tata:titi" {
		t.Log("Error TestFullPathImageFromDocker")
		t.Fail()
	}
}

func TestFullPathImageFromDockerWithTag(t *testing.T) {
	if getPatchedImageUrl("docker.io/toto/tata", "example.com") != "example.com/toto/tata:latest" {
		t.Log("Error TestFullPathImageFromDockerWithTag")
		t.Fail()
	}
}

func TestFullPathImageLibraryFromDocker(t *testing.T) {
	if getPatchedImageUrl("docker.io/busybox:titi", "example.com") != "example.com/library/busybox:titi" {
		t.Log("Error TestFullPathImageLibraryFromDocker")
		t.Fail()
	}
}

func TestFullPathImageLibraryFromDockerWithTag(t *testing.T) {
	if getPatchedImageUrl("docker.io/busybox", "example.com") != "example.com/library/busybox:latest" {
		t.Log("Error TestFullPathImageLibraryFromDockerWithTag")
		t.Fail()
	}
}

func TestForeignImage(t *testing.T) {
	if getPatchedImageUrl("gcr.io/busybox:titi", "example.com") != "gcr.io/busybox:titi" {
		t.Log("Error TestForeignImage")
		t.Fail()
	}
}

func TestForeignFullPathImage(t *testing.T) {
	if getPatchedImageUrl("gcr.io/toto/tata:titi", "example.com") != "gcr.io/toto/tata:titi" {
		t.Log("Error TestForeignFullPathImage")
		t.Fail()
	}
}

func TestRegistryImage(t *testing.T) {
	if getPatchedImageUrl("registry:titi", "example.com") != "docker.io/library/registry:titi" {
		t.Log("Error TestRegistryImage")
		t.Fail()
	}
}

func TestRegistryImageWithTag(t *testing.T) {
	if getPatchedImageUrl("registry:tag", "example.com") != "docker.io/library/registry:tag" {
		t.Log("Error TestRegistryImageWithTag")
		t.Fail()
	}
}

func TestRegistryPathImage(t *testing.T) {
	if getPatchedImageUrl("library/registry:titi", "example.com") != "docker.io/library/registry:titi" {
		t.Log("Error TestRegistryPathImage")
		t.Fail()
	}
}

func TestRegistryFullPathImage(t *testing.T) {
	if getPatchedImageUrl("docker.io/library/registry", "example.com") != "docker.io/library/registry:latest" {
		t.Log("Error TestRegistryFullPathImage")
		t.Fail()
	}
}
