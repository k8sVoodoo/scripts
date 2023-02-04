import docker
import os

#Need the below environment variables to be set
HARBOR_PROJECT="<projectname>"
HARBOR_URL="<harbor.url.example.com>"
HARBOR_USERNAME="<username>"
HARBOR_PASSWORD="<password>"

# Docker client
client = docker.from_env()

# Login to docker registry
client.login(
      username=HARBOR_USERNAME,
      password=HARBOR_PASSWORD,
      registry=HARBOR_URL
)

# Get images from docker client
images = client.images.list()

for img in images:
    # Tag images in the format required by Harbor.
    tag = f"{HARBOR_URL}/{HARBOR_PROJECT}/{img.tags[0]}"
    # Tag the image
    client.images.tag(image=img.id, tag=tag)
    # Push the image to the Harbor registry.
    client.images.push(image=tag, tag=tag)

# Log out of docker registry
client.logout(HARBOR_URL)
