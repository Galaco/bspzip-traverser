# bspzip-traverser
Generate a bspzip filelist.txt from a directory for packing into a Source Engine .bsp

Use your working directory however you like, Traverser can build a better and smaller filelist.

#####Provided executables are built with gox. You can either trust these builds, or build from source yourself (see below).


# Usage

### From command line
#### Simple
Standard usage is very easy. There are 2 mandatory flags to specify; `target` your directory to walk through, and `output` your resultant filelist.txt.
```bash
bspzip-traverser.exe -target="C:/mapname/mount/" -output="C:/mapname/mapname-filelist.txt"
```
This command will generate a `mapname-filelist.txt` from the entire stucture contained within `C:/mapname/mount/` in the directory `C:/mapname/`

#### Strict mode
There is a single optional flag that you may find very useful, that can be added like so:
```bash
bspzip-traverser.exe -target="C:/mapname/mount/" -output="C:/mapname/mapname-filelist.txt" -strict
```
The `-strict` flag is great for ensuring you definitely don't pack misplaced or source files. 

For example, you keep your material source images (e.g. `.tga`) in the same directory as your converted `.vtf`, or you misplaced it accidentally. Strict mode can automatically skip all filetypes Source Engine doesn't expect in a given top level directory. This would allow you to keep whatever you want in your mount directory. This can help you reduce your final map filesize without having to worry about decluttering your working directory.

Strict mode supports all known valid extensions per directory, so it will keep `vtf`'s in `/materials`, but will ignore them in `/models`

NOTE: Strict will allow all extensions in the top level directory, and any files that are not part of an expected Source directory

### Adding to Hammer as an automated build step
@TODO

### Building from source
Assuming you have golang installed, simply run go build in the project directory.