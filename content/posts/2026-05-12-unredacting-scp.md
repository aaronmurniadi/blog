---
layout: post
title: 'Using AI to unredact SCP Foundation archives'
date: 2026-05-12
last_modified_date: 2026-05-12
---

# Using AI to unredact SCP Foundation archives

I've spent way too much time on the [SCP Foundation
wiki](https://scp-wiki.wikidot.com/). If you've been there, you know the deal:
official-looking reports that suddenly hit you with a big black bar. Those
[REDACTED] boxes are basically what makes the site work. They're meant to be
scary because they hide the details and let you imagine the worst.

But lately, I've started seeing those black boxes as a bit of a challenge.

In these stories, the redactions hide secrets from the reader. When a file says
someone "ate the heart of [DATA EXPUNGED]," you aren't just looking at a
censored line. You're looking at a specific gap in a sentence. I started
wondering if I could use an AI to figure out what was actually supposed to be
there.

I've been trying a simple way to do this. I don't look at the black bar itself.
Instead, I look at the words around it. Even when a name or an object is hidden,
the way the sentence is written gives away what kind of thing is missing. There
is many ways to look at this, but I usually keep it simple.

I usually look at how the redacted part works in the sentence:

- Who is doing the action? (Example: "The object made [REDACTED] explode.")
- What is the action happening to? (Example: "Doctors tried to cut out
[REDACTED].")
- Is it something that can think? (Example: "The subject said [REDACTED]
whispered to him.")

Once I have those clues, I ask an AI to look at other SCP stories and give me
its best guess. I'm not asking it to just guess randomly; I'm asking it to find
what fits the vibe of the story best.

I tried this with a story about a weird inkwell. The sentence was: "The subject
started to [REDACTED] like they were reading their own obituary."

The AI gave me a few options:

1. Weeping or crying: The most normal human reaction.
2. Writing: This fits the inkwell and obituary theme.
3. Dissolving: A common SCP move where the person turns into ink.

The cool part is that sometimes the AI comes up with something even creepier
than what I thought of. I'm not trying to ruin the mystery. For me, the unknown
is the best part of SCP. But doing this helps me see how the writers build these
scares by leaving specific holes for me to fill.

If you want to try it, just copy an SCP article from the wiki and paste it   
to this [SCP Redaction Analysis Gemini
Gem](https://gemini.google.com/gem/1fXAYUKnnYpydzMEV78JtegbBdH-gEBGe?usp=sharing)

Note: the Deep Research feature is enabled by default and can take some time. If
you want faster results, you can turn it off.

It's a fun way to spend an afternoon. 

Or, do you think the mystery is better left alone?