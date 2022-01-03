{
  description = "A very basic flake";

  outputs = { self, nixpkgs }: {

    packages.x86_64-linux =
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      vendorSha256 = "sha256-rHaqxAb27amS36msMo3Ry70UWzTbuK/jT4HLKkeDP4Y=";
    in rec {
      patch-docker-image-name = pkgs.buildGoModule {
        CGO_ENABLED = "0";
        pname = "patch-docker-image-name";
        version = "0.1.0";

        src = ./.;

        inherit vendorSha256;
        subPackages = [ "cmd/patch-docker-image-name" ];
      };

      normalize-docker-image-name = pkgs.buildGoModule {
        CGO_ENABLED = "0";
        pname = "normalize-docker-image-name";
        version = "0.1.0";

        src = ./.;

        inherit vendorSha256;
        subPackages = [ "cmd/normalize-docker-image-name" ];
      };

      k8s-proxy-image-swapper = pkgs.buildGoModule {
        CGO_ENABLED = "0";

        pname = "k8s-proxy-image-swapper";
        version = "0.3.1";

        src = ./.;

        inherit vendorSha256;
        subPackages = [ "." ];
      };

      oci-k8s-proxy-image-swapper = pkgs.dockerTools.buildLayeredImage {
        name = "oci-k8s-proxy-image-swapper";
        contents = [ k8s-proxy-image-swapper ];
        config = {
          Entrypoint = [ "${k8s-proxy-image-swapper}/bin/k8s-proxy-image-swapper" ];
          User = "1000:1000";
        };
      };
    };

    defaultPackage.x86_64-linux = self.packages.x86_64-linux.k8s-proxy-image-swapper;

  };
}
