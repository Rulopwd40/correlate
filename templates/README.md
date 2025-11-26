# Correlate Templates Guide

Templates are JSON configuration files that define how to build, update, and manage dependencies for different project types.

## Template Structure

```json
{
  "name": "template-name",
  "description": "Human-readable description",
  "variables": {
    "projectIdentifier": "pattern for project ID",
    "version": "pattern for version"
  },
  "detect": {
    "manifest": "manifest-file-name",
    "searchPattern": "pattern to find dependency"
  },
  "steps": [
    {
      "name": "Step name",
      "type": "command|script",
      "cmd": "single command",
      "script": ["multiple", "commands"],
      "workdir": "{{variable}}",
      "variables": {
        "stepVar": "value"
      }
    }
  ]
}
```

## Variable System

### Global Variables

Available in all steps:
- `{{identifier}}` - The dependency identifier
- `{{sourceDir}}` - Directory of the source project
- `{{targetDir}}` - Directory of the dependency project

### Custom Variables

You can add custom variables in the template or at step level:

```json
{
  "steps": [
    {
      "name": "Build",
      "cmd": "ng build @{{scope}}/{{identifier}}",
      "workdir": "{{sourceDir}}",
      "variables": {
        "scope": "mycompany"
      }
    }
  ]
}
```

### Variable Resolution

Variables are resolved using `{{variableName}}` syntax:

```json
"cmd": "npm install {{sourceDir}}/dist/{{identifier}}.tgz"
```

## Step Types

### Command Type (Default)

Single command execution:

```json
{
  "name": "Build project",
  "cmd": "mvn clean install",
  "workdir": "{{sourceDir}}"
}
```

### Script Type

Multiple commands chained together:

```json
{
  "name": "Build and pack",
  "type": "script",
  "script": [
    "cd {{sourceDir}}/dist",
    "npm pack",
    "mv *.tgz {{targetDir}}"
  ],
  "workdir": "{{sourceDir}}"
}
```

## Example Templates

### Java Maven

```json
{
  "name": "java-maven",
  "description": "Template for Java projects using Maven",
  "variables": {
    "projectIdentifier": "<artifactId>{{projectIdentifier}}</artifactId>",
    "version": "<version>{{version}}</version>"
  },
  "detect": {
    "manifest": "pom.xml",
    "searchPattern": "<artifactId>{{identifier}}</artifactId>"
  },
  "steps": [
    {
      "name": "Build artifact",
      "cmd": "mvn clean install",
      "workdir": "{{sourceDir}}"
    },
    {
      "name": "Update pom versions",
      "cmd": "correlate replace {{identifier}}",
      "workdir": "{{targetDir}}"
    },
    {
      "name": "Recompile dependent projects",
      "cmd": "mvn clean compile",
      "workdir": "{{targetDir}}"
    }
  ]
}
```

### Angular Library (Local Tarball)

```json
{
  "name": "angular-library-local",
  "description": "Angular library with local tarball deployment",
  "variables": {
    "projectIdentifier": "\"name\": \"@{{scope}}/{{projectIdentifier}}\"",
    "version": "\"version\": \"{{version}}\""
  },
  "detect": {
    "manifest": "package.json",
    "searchPattern": "\"@{{scope}}/{{identifier}}\""
  },
  "steps": [
    {
      "name": "Update library version",
      "cmd": "npm version {{newVersion}} --no-git-tag-version",
      "workdir": "{{sourceDir}}",
      "variables": {
        "newVersion": "1.0.0-SN-2"
      }
    },
    {
      "name": "Build Angular library",
      "cmd": "ng build @{{scope}}/{{identifier}}",
      "workdir": "{{sourceDir}}"
    },
    {
      "name": "Pack library to tarball",
      "type": "script",
      "script": [
        "cd {{sourceDir}}/dist/{{scope}}/{{identifier}}",
        "npm pack"
      ],
      "workdir": "{{sourceDir}}"
    },
    {
      "name": "Install tarball in consumer",
      "cmd": "npm install {{sourceDir}}/dist/{{scope}}/{{identifier}}/{{scope}}-{{identifier}}-{{newVersion}}.tgz",
      "workdir": "{{targetDir}}"
    },
    {
      "name": "Install dependencies",
      "cmd": "npm install",
      "workdir": "{{targetDir}}"
    }
  ]
}
```

### Python Pip

```json
{
  "name": "python-pip",
  "description": "Template for Python projects using pip",
  "variables": {
    "projectIdentifier": "name=\"{{projectIdentifier}}\"",
    "version": "version=\"{{version}}\""
  },
  "detect": {
    "manifest": "setup.py",
    "searchPattern": "name=\"{{identifier}}\""
  },
  "steps": [
    {
      "name": "Install package",
      "cmd": "pip install -e .",
      "workdir": "{{sourceDir}}"
    },
    {
      "name": "Build distribution",
      "cmd": "python setup.py sdist bdist_wheel",
      "workdir": "{{sourceDir}}"
    },
    {
      "name": "Update requirements",
      "type": "script",
      "script": [
        "cd {{targetDir}}",
        "pip install --upgrade {{identifier}}"
      ],
      "workdir": "{{targetDir}}"
    }
  ]
}
```

## Usage

### Initialize Project

```bash
correlate init my-library java-maven
```

### Link Dependencies

```bash
correlate link my-library /path/to/consumer-project
```

### Update Dependencies

```bash
correlate update my-library
```

## Creating Custom Templates

1. Create a new `.json` file in the `templates/` directory
2. Define the manifest file name and search pattern
3. Add build/update steps with appropriate commands
4. Use variables for flexibility
5. Test with your project

## Tips

- Use `{{sourceDir}}` for the library being built
- Use `{{targetDir}}` for projects consuming the library
- Chain commands with `script` type for complex workflows
- Add step-level variables for per-step configuration
- Keep commands cross-platform when possible
