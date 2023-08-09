.PHONY: all
all: terraform.gif terraform-small.gif tf.gif tf-small.gif

%.gif: %.yml
	terminalizer render $< -o $@

%-small.gif: %.gif
	convert -verbose $< -resize 30% $@
