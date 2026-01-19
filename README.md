# Documentation Driven Development (DDD) Template

A template repository for building software projects where design decisions come first, captured as Design Records before (or during) implementation. Built for collaboration between humans and AI agents.

## What is Documentation Driven Development?

Documentation Driven Development (DDD) is an approach where design decisions are captured in documentation before (or during) implementation. While the term was originally used in 2004 academic research for real-time systems using formal specifications (Luqi et al.), this template takes a different approach: capturing design decisions in human-readable Design Records for modern software development with AI agent collaboration.

In practice, this means you:

1. Think first, code later - Make explicit design decisions before implementation
2. Document the "why" - Capture not just what you built, but why you chose that approach
3. Create Design Records (DRs) - Formal documents for each significant technical decision
4. Build institutional memory - Future developers (including AI agents) understand the reasoning behind choices

Unlike traditional "documentation after the fact," DDD makes documentation a first-class part of the design process.

## Why Use This Approach?

- **Better decisions** - Discussing alternatives and trade-offs upfront leads to more thoughtful choices
- **AI agent continuity** - Agents can pick up where previous sessions left off by reading DRs
- **Reduced churn** - Documented decisions prevent revisiting the same questions repeatedly
- **Team alignment** - Everyone understands not just what was built, but why
- **Knowledge preservation** - Critical context survives beyond the original developers

## Quick Start

1. **Clone or copy this template:**

   ```bash
   git clone <this-repo> my-new-project
   cd my-new-project
   rm -rf .git
   git init
   ```

2. **Update AGENTS.md:**
   - Replace `[Project Name]` with your project name
   - Add a brief description
   - Create your first project document when ready

3. **Create your first project:**

   ```bash
   # Create a project document for your initial work
   cp .ai/projects/p-writing-guide.md .ai/projects/p-001-initial-setup.md
   # Edit p-001 with your project goals and scope
   ```

4. **Set active project in AGENTS.md:**

   ```markdown
   Active Project: .ai/projects/p-001-initial-setup.md
   ```

5. **Start working with an AI agent:**
   - Agent reads AGENTS.md to find the active project
   - Agent reads the active project document
   - Agent follows task documents (design-phase.md or implementation-phase.md) as needed
   - Agent creates Design Records for significant decisions

## How It Works

### Projects Define What to Build

Create project documents (p-001, p-002, etc.) in `.ai/projects/` that define:

- Goals and scope
- Success criteria
- Deliverables
- Status tracking

See `.ai/projects/p-writing-guide.md` for how to structure projects.

### Design Records Document Why

As you work, create Design Records in `.ai/design/design-records/` for significant decisions:

- Architecture choices
- Technology selections
- Trade-off decisions
- Algorithm designs

See `.ai/design/dr-writing-guide.md` for how to write DRs.

### Task Documents Guide Agent Work

Task documents in `.ai/tasks/` provide focused guidance:

- `design-phase.md` - For making and documenting design decisions
- `implementation-phase.md` - For coding according to DRs

Reference these in your project documents or AGENTS.md to guide agent behavior.

### Workflow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. CLONE TEMPLATE                                               │
│    - Template includes p-001-bootstrap-template.md              │
│    - Update AGENTS.md with your project name/description        │
│    - Active project already set to p-001                        │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. DESIGN PHASE (.ai/tasks/design-phase.md)                    │
│    Agent reads:                                                 │
│    - AGENTS.md → finds active project (p-001)                   │
│    - p-001 → understands goals                                  │
│    - Existing DRs → learns previous decisions                   │
│                                                                 │
│    Agent activities:                                            │
│    - Ask design questions, explore alternatives                 │
│    - Create Design Records (DR-NNN-*.md)                        │
│    - Update .ai/design/design-records/README.md                │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. IMPLEMENTATION PHASE (.ai/tasks/implementation-phase.md)    │
│    Agent reads:                                                 │
│    - All DRs → understand design decisions                      │
│    - Active project → know what to build                        │
│                                                                 │
│    Agent activities:                                            │
│    - Implement according to DRs                                 │
│    - Add DR references in code comments                         │
│    - Write tests, create new DRs if gaps found                  │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. CODE REVIEW (.ai/tasks/code-review.md)                      │
│    - Review correctness, design, testing                        │
│    - Verify DR alignment                                        │
│    - Create rectification project if issues found               │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. PROJECT COMPLETION                                           │
│    - Verify success criteria met                                │
│    - Status: Pending → Active → Completed                       │
│    - Move to .ai/projects/completed/                           │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. NEXT PROJECT                                                 │
│    - Create p-002-*.md (or p-NNN)                               │
│    - Set as active in AGENTS.md                                 │
│    - Repeat cycle                                               │
└─────────────────────────────────────────────────────────────────┘
```

### Document Relationships

```
AGENTS.md (entry point)
    │
    ├─→ Active Project: .ai/projects/p-NNN-*.md
    │       │
    │       ├─→ References: .ai/tasks/design-phase.md
    │       │                   │
    │       │                   └─→ Creates: DR-NNN-*.md
    │       │
    │       └─→ References: .ai/tasks/implementation-phase.md
    │                   │
    │                   └─→ Reads DRs → Implements → References DRs in code
    │
    └─→ Quick Reference (populated as project develops)
```

## Repository Structure

```
.
├── README.md                     # This file - DDD overview
├── AGENTS.md                     # Agent entry point and active project
├── LICENSE                       # MIT license
├── .ai/                          # AI agent working files
│   ├── README.md                # .ai directory overview
│   ├── workflow.md              # Feature development workflow
│   ├── design/
│   │   ├── README.md            # Design docs overview
│   │   ├── dr-writing-guide.md  # How to write Design Records
│   │   └── design-records/
│   │       ├── README.md        # DR index
│   │       └── superseded/      # Superseded DRs (kept for history)
│   ├── projects/
│   │   ├── README.md            # Project index
│   │   ├── p-writing-guide.md   # How to write project documents
│   │   └── completed/           # Completed projects
│   └── tasks/
│       ├── README.md            # Task documents overview
│       ├── design-phase.md      # Guide for design work
│       ├── implementation-phase.md # Guide for implementation work
│       └── code-review.md       # Guide for code review
├── docs/                         # Human-facing documentation
│   └── thoughts.md              # Brainstorming and ideas
└── .gitignore
```

## Key Concepts

### Design Records (DRs)

Design Records are structured documents that capture:

- **Problem** - What constraint or issue drove this decision?
- **Decision** - What did we choose?
- **Why** - Core reasoning behind the choice
- **Trade-offs** - What we accept vs. what we gain
- **Alternatives** - What else we considered and why we rejected it

DRs are numbered sequentially (dr-001, dr-002, etc.) and include a category in the filename:

- `dr-001-config-file-format.md`
- `dr-002-api-authentication-strategy.md`

See `.ai/design/dr-writing-guide.md` for complete guidelines.

### Project Documents

Project documents define work packages with:

- Goals and objectives
- Scope (in/out)
- Success criteria
- Deliverables
- Status tracking

Projects are numbered sequentially (p-001, p-002, etc.):

- `p-001-initial-architecture.md`
- `p-002-user-authentication.md`

See `.ai/projects/p-writing-guide.md` for complete guidelines.

### Task Documents

Task documents in `.ai/tasks/` guide agents through specific types of work:

- **design-phase.md** - Making design decisions and creating DRs
- **implementation-phase.md** - Coding according to design decisions
- **code-review.md** - Reviewing code for quality and DR alignment

Reference these in project documents or AGENTS.md when you want focused agent behavior. They provide detailed workflows and instructions for agents working on different phases of development.

## Working With AI Agents

This template is optimized for AI agent collaboration:

1. **AGENTS.md is the entry point** - Points to the active project
2. **Active project provides context** - What we're working on now
3. **DRs provide continuity** - Why decisions were made
4. **Task documents guide behavior** - How to approach the work

Example agent prompt:

```
Read AGENTS.md, then read the active project document.

We're working on design decisions for the authentication system.
Follow .ai/tasks/design-phase.md to help explore options and create
Design Records.
```

## Documentation

### Writing Guides

- **DR Writing Guide**: `.ai/design/dr-writing-guide.md`
- **Project Writing Guide**: `.ai/projects/p-writing-guide.md`

### Workflow Guide

- **Feature Development Workflow**: `.ai/workflow.md` - Interactive step-by-step process for developing features collaboratively

### Task Guides

- **Design Phase**: `.ai/tasks/design-phase.md`
- **Implementation Phase**: `.ai/tasks/implementation-phase.md`

### Indexes

- **Design Records Index**: `.ai/design/design-records/README.md`
- **Projects Index**: `.ai/projects/README.md`

## Philosophy

**Documentation Driven Development is about thinking, not bureaucracy.**

Create DRs for decisions that:

- Have multiple valid approaches
- Involve trade-offs
- Will be hard to change later
- Future developers will wonder about

Don't create DRs for:

- Trivial implementation details
- Decisions with only one reasonable option
- Standard practices

Create projects for:

- Focused efforts with clear deliverables
- Work packages that can be scoped and completed
- Multi-step initiatives requiring planning

Don't create projects for:

- Single bug fixes
- Trivial changes
- Day-to-day maintenance

The goal is to capture meaningful decisions and organize substantial work, not document everything.

## License

MIT License - See [LICENSE](LICENSE) for details.

When using this template for your own project, you may delete or replace the LICENSE file. The template's license only governs the template files themselves, not projects created from it.

## Contributing

This is a template repository. If you have improvements to the template itself, feel free to fork and adapt to your needs.

---

**Ready to build something thoughtfully?** Clone this template and start your DDD journey.
