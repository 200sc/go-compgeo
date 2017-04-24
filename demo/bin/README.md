In each zip, there are four items:

  An executable to run, which starts the demo.

  A "oak.config" file, which provides configuraton settings to the engine.
  
  A "asssets" folder, which inside has a single font the engine uses to render with.
  
  And a "data" folder, which has a number of OFF files to load and test the program with.


The build.go script here automatically cross-compiles the demo package to
several os and architecture combinations. The result of running this file
can be combined with the assets, oak.config, and data files in the demo directory
to produce the contents of each zip file.