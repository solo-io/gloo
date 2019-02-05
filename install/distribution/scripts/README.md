# GlooE installation setup scripts
These two scripts are meant to support the user in adding `glooctl` to their PATH:
- setup.sh: adds the correct Unix-based executable to `~/.gloo/bin` and creates a symlink to it in `/usr/local/bin`
- setup.bat: adds the `glooctl` Windows executable file to the `%USERPROFILE%/.gloo/bin` directory. Windows does not have an 
equivalent for `/usr/local/bin`, so the only way to add the file to the `PATH` is to add the folder to the environment 
variable itself. Doing this in a script is non-trivial, as there is a risk of corrupting the `PATH`. We therefore just 
print out instructions for the users to do it themselves.