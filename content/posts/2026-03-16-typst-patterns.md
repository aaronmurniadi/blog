---
title: "How I Found My Typst Patterns (and Why the Community Needs More)"
date: 2026-03-16
layout: post
---

# How I Found My Typst Patterns (and Why the Community Needs More)

> [Check out my Typst templates here](/typesettings)!

While I worked on the [Harvard Law Review–style journal template](/typesettings#journal-template), I rearranged the file one more time. I moved one big `meta` dict to the top, with everything the issue needs to know about itself. Then I put the styling and layout after it. Opening the file then meant reading the document first and the machinery second.

```typst
#let meta = (
  journal: "Typst Type Review",
  year: "1958", 
  month: "February",
  volume: "71",
  number: "4",
  title: "The Foundations of Visual Language: The Law of Typography",
  author: "Aaron P. Murniadi",
  abstract: [...],
  body: [...],
)
```

That split, data up front and presentation below, felt obvious in hindsight. But it took a few failed layouts to get there. It also made the larger problem visible. Typst gives you a lot of rope. But it gives you almost no shared vocabulary for how to organize a real project.

## Docs teach syntax, not structure

The [official Typst documentation](https://typst.app/docs/) is strong on mechanics. When I need the exact behavior of `show` rules or page setup, I still reach for it first. It does not try to prescribe how to lay out a fifty-page report, a book, or a journal issue. It leaves that layout to you, whether on disk or in memory. Ask "how should I structure this?" in the community, and you will get several incompatible answers. There is no "house style" for project shape the way many ecosystems eventually develop.

So this post is not a complaint about the docs. It is a note that we are still missing a layer above them. We need conventions for document architecture, and honest writeups of what people actually do when the tutorial ends.

## What I reach for, depending on the job

### Plain `#set` and content

For drafts, short letters, or anything disposable, I skip abstraction entirely:

```typst
#set page(...)
#set text(...)
#set par(...)

#align(center)[SUPREME COURT OF THE UNITED STATES]

// Direct content without complex structure
```

No template wrapper, no shared config object. Just page and paragraph rules, then text. That is often the right amount of structure.

### One file, top-down

For [my CV](/typesettings#cv) and other one-off documents, I still use a single file. It holds the imports, a few helpers, then the body. Everything reads in order.

```typst
#import "@preview/droplet:0.3.1": dropcap

#set document(title: "Aaron P. Murniadi's CV")

#let section-block(title, content) = [
  #text(size: 1.25em, style: "italic", title)
  #block(inset: (left: 2em))[#content]
  #v(0.5em)
]

// ... more function definitions ...

#header(name: [Aaron P. Murniadi], contact: [...])
```

It is fast to write and easy to read the first time. It falls apart when you need variants or long-term maintenance. I tried multiple CV versions this way, and I ended up duplicating chunks of logic. For anything with a future, I move on.

### Splitting style from content

The book-length ["Maid of Orleans"](/typesettings#maid-of-orleans) project taught me that one giant file does not scale. The layout I use now looks like this:

```
maid_of_orleans/
├── maid_of_orleans.typ       # Main content
├── maid_of_orleans_style.typ # Style definitions
└── main.typ                  # Alternative layout
```

The style module owns the template and helpers:

```typst
#import "@preview/droplet:0.3.1": dropcap
#import "@preview/typearea:0.2.0": typearea

#let template(body) = {
  // Page setup and styling
  show: typearea.with(...)
  set text(...)
  set par(...)
  // Custom styling rules
  body
}

#let framed-image(img-path, cap) = { ... }
#let dropped(first, rest) = { ... }
```

The manuscript imports what it needs and stays mostly prose:

```typst
#import "maid_of_orleans_style.typ": dropped, framed-image, template, typearea
#show: template

// Cover and content
```

The tradeoff is navigation. You jump between files and need to remember which symbol lives where. For book-sized work, that cost has been worth it.

### A single `template(...)` wrapper

For submissions where someone else's spec fixes the margins and fonts, I use one function. It wraps the whole document and applies all the rules in one place:

```typst
#let template(
  title: none,
  abstract: none,
  authors: (),
  date: "©2025",
  doc,
) = {
  set page(...)
  set text(...)
  set par(...)
  // All styling rules
  doc
}

#show: doc => template(
  title: [Nine Things I Learned in Ninety Years],
  authors: ([Edward Packard],),
  abstract: [...],
  doc,
)
```

That pattern is easy to overuse. I spent evenings tuning the wrapper instead of writing. I reserve it for cases where the format really is non-negotiable.

### Configuration objects (including "meta first")

When the same codebase must serve more than one layout, I push the fields into one structure. This works when the document is a bundle of fields such as title, authors, abstract, and body. The rest of the file then consumes that structure.

Example output for a two-column article: [here](/typesettings#two-column-article).

```typst
#let config(
  column: 2,
  size: 10pt,
  font: "Libertinus Serif",
  paper: "a4",
  title: none,
  authors: (),
  abstract: [],
  doc,
) = {
  // Apply configuration
  set page(paper: paper, ...)
  set text(font: font, size: size, ...)
  // ... rest of styling
}

#show: document => config(
  column: 2,
  title: [Article Title],
  authors: (...),
  abstract: [...],
  document,
)
```

The journal template is the same idea with a dict instead of function arguments. A `meta` object holds the facts, then `#set`, `show`, and `#meta.body` wire it up.

```typst
#let meta = (
  institution: "Typst University",
  journal: "Typst Type Review", 
  year: "1958",
  title: "The Foundations of Visual Language",
  author: "Aaron P. Murniadi",
  abstract: [...],
  body: [...]
)

// Styling and setup using meta object
#set text(...)
#set page(...)

// Document content
#meta.body
```

Edits to journal name, volume, or author string stay localized. Styling can change without touching the content block.

## How I choose

I do not have one pattern for everything. A plain letter stays flat. A one-off can live in a single file until I need variants. Book-length work gets a style module. Conference-style specs get a wrapping template. Journals and anything with multiple outputs lean on a config object or a `meta` dict. Who will maintain it matters as much as document length. It could be just me on a deadline, or other people over months. I also assume the first layout will be wrong sometimes. Refactoring once I understand the content has been routine, not a mistake.

## Toward shared patterns

The ecosystem would benefit from more public examples of "how we structured this" alongside "how the syntax works". This matters most for large or regulated documents. I explained what works on my machine. I would like to read the same kind of post from others. I want the layouts that held up, the layouts that did not, and why.

What organizational habits have you settled on in Typst? What broke the first time you tried it?
