---
title: "How I Found My Typst Patterns (and Why the Community Needs More)"
date: 2026-03-16
layout: post
---

# How I Found My Typst Patterns (and Why the Community Needs More)

> [Check out my Typst templates here](/typesettings)!

While I was working on the [Harvard Law Review–style journal template](/typesettings#journal-template), I rearranged the file one more time: one big `meta` dict for everything the issue needs to know about itself, then the styling and layout after that. Opening the file suddenly meant reading the document first and the machinery second.

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

That split—data up front, presentation below—felt obvious in hindsight, but it took a few failed layouts to get there. It also made the larger problem visible: Typst gives you a lot of rope, and almost no shared vocabulary for how to organize a real project.

## Docs teach syntax, not structure

The [official Typst documentation](https://typst.app/docs/) is strong on mechanics. When I need the exact behavior of `show` rules or page setup, I still reach for it first. What it does not try to do is prescribe how a fifty-page report, a book, or a journal issue should be laid out on disk or in memory. Ask "how should I structure this?" in the community and you will get several incompatible answers; there is no "house style" for project shape the way many ecosystems eventually develop.

So this post is not a complaint about the docs. It is a note that we are still missing a layer above them: conventions for document architecture, and honest writeups of what people actually do when the tutorial ends.

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

No template wrapper, no shared config object—just page and paragraph rules, then text. That is often the right amount of structure.

### One file, top-down

For [my CV](/typesettings#cv) and other one-off documents, I still use a single file with imports, a few helpers, then the body. Everything reads in order.

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

It is fast to write and easy to read the first time. It falls apart when you need variants or long-term maintenance: I tried multiple CV versions this way and ended up duplicating chunks of logic. For anything with a future, I move on.

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

The tradeoff is navigation: you jump between files and need to remember which symbol lives where. For book-sized work, that cost has been worth it.

### A single `template(...)` wrapper

For submissions where margins and fonts are fixed by someone else’s spec, I use a function that wraps the whole document and applies all the rules in one place:

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

That pattern is easy to overuse. I have spent evenings tuning the wrapper instead of writing. I reserve it for cases where the format really is non-negotiable.

### Configuration objects (including "meta first")

When the same codebase might serve more than one layout, or when the document is really a bundle of fields (title, authors, abstract, body), I push those fields into a single structure and let the rest of the file consume it.

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

The journal template is the same idea with a dict instead of function arguments: one `meta` object holds the facts, then `#set` / `show` and the final `#meta.body` wire it up.

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

Edits to journal name, volume, or author string stay localized; styling can change without touching the content block.

## How I choose

I do not have one pattern for everything. A plain letter stays flat; a one-off can live in a single file until I need variants; book-length work gets a style module; conference-style specs get a wrapping template; journals and anything with multiple outputs lean on a config object or a `meta` dict. Who will maintain it—just me on a deadline, or other people over months—matters as much as document length. I also assume the first layout will be wrong sometimes: refactoring once I understand the content has been routine, not a mistake.

## Toward shared patterns

The ecosystem would benefit from more public examples of "how we structured this" alongside "how the syntax works"—especially for large or regulated documents. I have written up what works on my machine; I would like to read the same kind of post from others: layouts that held up, layouts that did not, and why.

What organizational habits have you settled on in Typst? What broke the first time you tried it?
