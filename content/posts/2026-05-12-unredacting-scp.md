---
layout: post
title: 'Using AI to unredact SCP Foundation archives'
date: 2026-05-12
last_modified_date: 2026-05-12
---

# Using AI to unredact SCP Foundation archives

I spent too much time on the [SCP Foundation
wiki](https://scp-wiki.wikidot.com/). If you know the site, you know the deal.
Official-looking reports suddenly hit you with a big black bar. Those
[REDACTED] boxes basically make the site work. They are meant to be scary,
because they hide the details and let you imagine the worst.

But lately I started to see those black boxes as a challenge.

In these stories, the redactions hide secrets from the reader. When a file says
someone "ate the heart of [DATA EXPUNGED]," you are not just looking at a
censored line. You are looking at a specific gap in a sentence. I wanted to
know if I could teach an AI to figure out what was actually supposed to be
there.

I tried a simple way to do this. I do not look at the black bar
itself. Instead, I look at the words around it. Even when a name or an object is
hidden, the sentence structure gives away what kind of thing is missing. There
are many ways to look at this, but I usually keep it simple.

I usually look at how the redacted part works in the sentence:

- Who is doing the action? (Example: "The object made [REDACTED] explode.")
- What is the action happening to? (Example: "Doctors tried to cut out
[REDACTED].")
- Is it something that can think? (Example: "The subject said [REDACTED]
whispered to him.")

Once I have those clues, I ask an AI to look at other SCP stories and give me
its best guess. I do not ask it to guess randomly. I ask it to find the option
that fits the vibe of the story best.

I tried this with a story about a weird inkwell. The sentence was: "The subject
started to [REDACTED] like they were reading their own obituary."

The AI gave me a few options:

1. Weeping or crying: The most normal human reaction.
2. Writing: This fits the inkwell and obituary theme.
3. Dissolving: A common SCP move where the person turns into ink.

The cool part is that sometimes the AI comes up with something even creepier
than what I thought of. I do not try to ruin the mystery. For me, the unknown
is the best part of SCP. But doing this helps me see how the writers build these
scares by leaving specific holes for me to fill.

If you want to try it, copy an SCP article from the wiki and paste it   
to this [SCP Redaction Analysis Gemini
Gem](https://gemini.google.com/gem/1fXAYUKnnYpydzMEV78JtegbBdH-gEBGe?usp=sharing)

Note: the Deep Research feature is enabled by default and can take some time. If
you want faster results, you can turn it off.

It is a fun way to spend an afternoon. 

Or, do you think the mystery is better left alone?