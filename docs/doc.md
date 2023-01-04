## Blip Template Engine

The blip template engine converts text files with blip commands into go files that output the template to the writer.

The blip template translator converts x.blip files to x.blip.go files.  This command is:

**blip**
```
blip --help
Blip Template Compiler
Blip Processing: Version: x.x.x
  -dir string
    	The source directory containing templates (default "./template")
  -help
    	Print help message
  -rebuild
    	rebuild all files
  -supportBranch string
    	Support branch name for include. (default "github.com/samlotti/blip/blipUtil")
  -watch
    	will watch the directory for file names/new files
```


# Installing blip.

go get ...
go install ..


# Blip supports 

* extending templates with placeholders for content
* including other templates 
* looping structures 
* conditionals
* variables to be passed into templates
* context variables
* text output
* embedded go code 
* go functions in templates

# Blip Commands
Commands are activated with the @ character.  In order to place an @ literally in the output it will need to be escapes as @@.  Note that @text will stop the @ escape requirements until @end is encountered.

In the blocks that follow, we will consider the output of a list of users in the output.
The user is defined as 
```go
type User struct {
	UId string
	Name string
	Active bool
}
var bob = User{Name="bob" ... }
var allUsers = []User{ bob, User2, .. UserN  }
 
```


#### @=     ... @
This is the basic replacement command.  It will render the content into the output stream.
<strong>Name: @= bob.Name@</strong>
Note that the code between @= and @ will be inserted into the output as go code and run during the template execution runtime.

The output will be escaped to make it safe based on the file type. The file type is determined by the file name. 
Ex: index.blip.html --> html file
Ex: amounts.blip.csv --> csv file (a implementation of IBlipEscaper would need to be registered to BlipUtil.Instance().AddEscaper)

For example if bobs name was <b>Bob</b>, then it will appear as <b>Bob</b> in the output.
This should be used for user generated content.

#### @==    ... @
The raw unescaped version of output. 
For example if bobs name was <b>Bob</b>, then it will appear bold in the output.
This should never be used for user generated content.

#### @int=    ... @  and @int64=    ... @
The value must be int or int64 types and will for formatted using iToA conversion.

#### @text  ... @end
Writes the content to the output.  This is only needed if there is some @ signs in the content and don't want to escape with @@.
Note that between commands in the template the content is written to the output.

#### @arg name type
This is how you can pass parameters to the template.  The @arg must at the root level and be the only command on the line.
In our example we can pass a user as follows:
@arg user *User
@arg allUsers []*User

The result of this will be the Render function requiring these parameters in the call.
```go
func UserListEntryRender( user *User, allUsers []*User, c context.Context, w io.Writer ) (terror error) {
	// ..
}
```

#### @context name type
These are variables that are passed via the context.  Within the template they are accessed directly as they get resolved at the start of the process.
If is important that the context has the entry otherwise the template will return a *runtime* error.

This is great for common objects such as the loggedIn user and so in.

#### @context name type = some initial value
In this case the initial value is used if the value is not in the context.

#### @// Single line comment
@// This comment will not be included in the generated code.


#### @* multiline comment  ..... til .. *@
@* This is a multi line comment
it will not be included in the generated code.

Can be used to comment out some template code.

*@

#### @include templateName arg1, arg2
#### @include package.templateName arg1, arg2   <-- If not in the same package as calling template
Calls the templateName render (generated) function passing in the arguments if any.

Example:
```
  @// Renders the user template for each user
  @for user in users
    <div>
    @include userPanel user
    </div>
  @end
``` 

#### @extend templateName arg1, arg2   ... @end
#### @extend package.templateName arg1, arg2   ... @end    <-- If not in the same package as calling template
Calls the render of the template. This is block command that can contain @content sections, the sections are rendered in the extended template.

The template being expected would have @yield commands that will output content from the content sections.

Example:

```html
    templateLayout
        <div>
            @yield topContent
        </div>
        <div>
            @yield midContent
        </div>
        <div>
            @yield bottom
        </div>
      
    child template
    @extend templateLayout
        @content bottom
            Number Users: @int= len(users) @
        @end
        @content top
            <title>Users</title>
        @end
        @content midContent
            @for user in users
            <div>
                @include userPanel user
            </div>
            @end
        @end

    @end
    

```

### @yield contentName
Renders the content from the caller

### @content contentName ... @end
Provides content for the extended template.

### @if @then @else
@if .. go expression
    template content
@else
    template content
@end

```
@if user.Active
	@include component.activeUser user
@else
	@include component.inActiveUser user
@end
```

# File names
Blip files are identified with the following patterns.

file.blip - This is a basic text blip file.  There will be no escaping of output.
@= and @== do not escape.

file.blip.html - Html file output, escaping for html tags.  @= escapes, @== does not escape

file.blip.{someOther} - Custom extension.
@== Will 

file.blip.html.go - The result file of the template engine.  This is a standard go file and should not be edited as it will be overridden.
Note that there is no blip version or date generated to make the file safe for source control as that only real changes will be detected.


# Escape strings
By default, blip support a null Escape for text (no escaping) and html for html files (html.EscapeString).
If you add your own type then you can also add an escaper.

To add an escaper follow these steps:
* Create a class that implements IBlipEscaper
* Add an instance to BlipUtil
* BlipUtil.Instance().AddEscaper("myType", myEscaper)
Note this should be configured on startup.

# Measure and monitoring
BlipUtil.Instance().SetMonitor( IBlipMonitor ) can be registered to have callbacks when a template has completed.

The default is no monitoring.  You can always monitor the call to the template render yourself.

IBlipMonitor will receive a call from each render including all include / extend calls.

# Generated files.

All generated go files will be in a subdirectory off the root subdirectory of the project.  The directory is called 'blipped'.
The templates are within the package name based on the source template package name.  Note that this will be flattened to one level.

Ex:

    web/template/components
    web/template/layout
    web/template/pages


Will generate go files in:


    blipped/components
    blipped/layout
    blipped/pages








