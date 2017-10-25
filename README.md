# gmb - The Go Media Bot
Automates monotone media handling when writing blog posts with a static site generator such as [hugo](https://github.com/gohugoio/hugo).

gmb creates directory structure, resizes images and create a stub blog post containing tags for all your images.

## Usage scenario
* You want to write a blog post containing lots of pictures using [hugo](https://github.com/gohugoio/hugo)
* You keep you original pictures in full resolution in a separate directory somewhere.
* You want to generate smaller versions of the pictures and output them to a suitable folder in your hugo source tree such as: "static/images/slug/"
* You want to generate a stub blog post such as "content/stub.md" containing the front matter and image tags for all your images (or audio/pdfs etc.).

## What you do
* You create a config file in the same folder as your originals.

config.yaml
``` yaml
base_dir: /home/username/hugosite/src # The path to your hugo site source.
slug: awesome-post                    # The slug you want your post to have.
img_limit: 1000                       # The resoultion limit as a square, defaults to 800x800 pixels.
```

then you run `gmb config.yaml`

gmb will then do the following:

* create /home/username/hugosite/src/static/images/awesome-post
* process all images in original folder and resize and save the output in the static/images/awesome-post folder. The process strips the exif information but rotates the image if needed.
* create /home/username/hugosite/src/content/awesome-post.md
* populate the post with frontmatter and image tags for all the processed images.
