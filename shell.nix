{pkgs ? import <nixpkgs> {}}:

let
  templ = builtins.getFlake("github:a-h/templ");
in
pkgs.mkShell {
	nativeBuildInputs = with pkgs; [
		go

		gopls
		delve
		go-tools

		air
		tailwindcss

		templ.packages.x86_64-linux.templ
	];
}

