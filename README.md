# Correlate

**Correlate** is a command-line interface (CLI) tool for **dependency management across projects**, based on template-driven pipelines. It allows you to automate build tasks, version updates, and custom steps across multiple projects.

---

## Key Features

- Initialize projects from **predefined templates** (Java/Maven, Angular/NPM, etc.)  
- Automate version and dependency updates  
- Template-driven task pipelines: build, test, compile, version replacement  
- Supports both local projects and remote templates from GitHub  
- Task-level logging for easier **debugging**  
- Simple CLI commands: `init`, `link`, `replace`, `update`

---

## Main Commands

| Command | Description |
|---------|-------------|
| `init <project> <template>` | Initializes a project using an existing template |
| `link [identifier] [fullPath]` | Links dependent projects together |
| `replace <identifier>` | Replaces versions or identifiers in configuration files |
| `update [identifier]` | Runs the task pipeline for a project or a specific reference |

---

## Typical Workflow

```bash
# Initialize a Java project using Maven
correlate init my-project java-maven

# Link dependent projects
correlate link

# Update identifiers/versions in the project
correlate replace my-project

# Run the full pipeline
correlate update
