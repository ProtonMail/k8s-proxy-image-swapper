{
  description = "A very basic flake";

  outputs = { self, nixpkgs }: {

    packages.x86_64-linux =
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };

    in rec {
      k8s-proxy-image-swapper = pkgs.buildGoModule {
        CGO_ENABLED = "0";

        pname = "k8s-proxy-image-swapper";
        version = "0.3.0";

        src = ./.;

        vendorSha256 = "sha256-lJ/DPH8XqFHYtIBAwB9BH/TFc5Q4t2tfdx8uNcky+ek=";
        subPackages = [ "." ];
      };

      oci-k8s-proxy-image-swapper = pkgs.dockerTools.buildLayeredImage {
        name = "oci-k8s-proxy-image-swapper";
        contents = [ k8s-proxy-image-swapper ];
        config.Cmd = [ "${k8s-proxy-image-swapper}/bin/k8s-proxy-image-swapper" ];
      };
    };

    defaultPackage.x86_64-linux = self.packages.x86_64-linux.k8s-proxy-image-swapper;

  };
}
