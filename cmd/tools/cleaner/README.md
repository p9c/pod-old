# Cleaner

Cleaner simply just reads in a source file, arranges the the source in the following order, grouping all root level declarations:

- types
- const
- var
- func

It simply sorts them lexicographically inside these groups. 

Constants are left as groups because of iota increments, but variables and types are all split into one per declaration.

The function sorting sorts by *, receiver type and function name, in that order, and being that ASCII puts the capital letters first, the exported methods are higher up and unexported lower down, within each of the classes in the file.

A key purpose of this little tool is to assist with splitting and merging source files within packages so as to reduce the amount of scrolling required, and so logically connected things are grouped together rather than hodge podge. 

Variable, type and constant declarations are often tempting to mix between functions near where the data / type is used, but again, this can make it difficult to get a clear and quick gauge of the contents of a source file.

I personally have a rule if I have to scroll through more than two screenfuls the source is probably getting too long anyway (and probably is a lot copypasta), and as I attempted to clean up the source in this project repository, hoping I would finish before the second coming, this tool would have saved me so much time.

Long source files and huge APIs are a maintenance nightmare. Hopefully this will help you avoid that, or more quickly deal with a mess you inherited from some former C++ programmer.

## Known issues

- Methods without a receiver name will disappear when you run cleaner over them. Give them names, even if you aren't using the receiver. Many of these 'omit if implicit' rules in Go complicate parsing. I am not sure why exactly this one type of function slips through but the issue may be resolved at some point. In VSCode the regex `func [(]\*` will find at least the pointers and `func (_ *` will change it so it doesn't break (work in progress, maybe I fix later but for now supervise it).

- Parenthesised var and type declarations are split (so they sort properly) but due to quirks in the `dst` library (a fork of go/ast that keeps whitespace and comments together with the code they refer to) it is not simple to break the group and properly format it. Manual removal of the braces will be required.
