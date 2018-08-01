# docker-sphinxsearch-alpine
Docker image for SphinxSearch based on the Alpine Linux 3.6 image

Make sure you map a config volume and generate the indexes (with indexer), or else the container won't start.

Please note that this image uses the Alpine edge/testing repository (as Sphinx is not yet in the main repo)

# ref
repo: https://github.com/marijngiesen/docker-sphinxsearch-alpine