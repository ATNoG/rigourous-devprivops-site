# DevPrivOps Visualizer

This repository holds the code for the tool's visualizer.

## Deployment

1. Install all dependencies (go gopls delve go-tools air tailwindcss templ) or run `nix-shell` to install them all.
2. Run `go build` to compile the visualizer
3. Configure the `.env` file according to the documentation in `.env.example`
4. Run `./devprivops-dashboard` to launch the visualizer

