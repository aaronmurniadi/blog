#set page(
  paper: "a6",
)

#set text(
  hyphenate: true,
  lang: "en",
  number-width: "proportional",
  costs: (hyphenation: 30%),
  size: 10pt,
  number-type: "old-style",
  font: "Iosevka NF",
)

#set par(
  first-line-indent: (amount: 2em, all: true),
  leading: 0.65em,
  spacing: 0.65em,
  justify: true,
  justification-limits: (
    tracking: (max: 0.1em, min: -0.015em),
  ),
)

#let img(source, page_break: true, overlay: none) = {
  if page_break {
    pagebreak()
  }
  context {
    let pg = counter(page).get().first()
    let alignment = none
    let dx = none
    let dy = none
    let image_width = none
    let image_height = none
    let image_fit = "cover"
    let bleed = 1mm

    if calc.odd(pg) {
      alignment = top + left
      dx = -page.margin.inside
      dy = -page.margin.top
      image_width = 100% + page.margin.outside + page.margin.inside + bleed
      image_height = 100% + page.margin.bottom + page.margin.top + bleed
    } else {
      alignment = top + right
      dx = page.margin.inside + bleed
      dy = -page.margin.top
      image_width = 100% + page.margin.outside + page.margin.inside + bleed
      image_height = 100% + page.margin.bottom + page.margin.top + bleed
    }

    place(
      alignment,
      dx: dx,
      dy: dy,
      stack(
        image(
          source,
          width: image_width,
          height: image_height,
          fit: image_fit,
        ),
      ),
    )

    if overlay != none {
      place(
        alignment,
        dx: dx,
        dy: dy,
        box(width: image_width, height: image_height, fill: overlay),
      )
    }
  }
  if page_break {
    pagebreak()
  }
}

#set page(margin: (top: 10mm, bottom: 10mm, inside: 10mm, outside: 10mm))

#[
  #set par(first-line-indent: 0em)
  #img(
    "../images/chair_fujifilm_c200_kodak_2383.jpg",
    page_break: false,
    overlay: gray.transparentize(70%),
  )

  #v(10fr)

  #text(
    size: 32pt,
    weight: "bold",
    fill: white,
  )[Numinous]

  #v(1fr)

  #text(fill: white)[ˈno͞omənəs]

  #v(1fr)

  #text(fill: white, style: "italic")[
    An adjective used to describe something that has a deep, mysterious, or
    spiritual quality. It often refers to an awe-inspiring, supernatural
    presence or atmosphere that surpasses human comprehension and appeals to a
    sense of the holy or transcendent.
  ]

  #v(1fr)
]

#pagebreak(to: "odd")

#set page(
  margin: (top: 10mm, bottom: 10mm, inside: 20mm, outside: 10mm),
  flipped: false,
)

#counter(page).update(1)

#set page(
  header: context {
    let pg = counter(page).get().first()
    set text(size: 0.9em, number-type: "lining")
    if pg == 1 {} else {}
  },

  footer: context {
    let pg = counter(page).get().first()
    set text(
      size: 1.2em,
      number-type: "lining",
      weight: "bold",
      fill: white,
    )
    let box-alignment = none
    let box-size = page.margin.bottom - 0.45mm
    let box-fill = black
    let box-size = page.margin.bottom - 0.45mm
    let box-align = center + horizon
    let box-rotate = 0deg
    let box-dx = none
    let box-dy = none
    let bleed = 1mm

    if calc.even(pg) {
      box-alignment = bottom + left
      box-dx = -page.margin.outside - bleed
      box-dy = page.margin.bottom - box-size
      box-align = center + horizon
      box-rotate = -90deg
    } else {
      box-alignment = bottom + right
      box-dx = page.margin.outside + bleed
      box-dy = page.margin.bottom - box-size
      box-align = center + horizon
      box-rotate = 90deg
    }

    place(
      box-alignment,
      dx: box-dx,
      dy: box-dy,
      block(
        fill: box-fill,
        width: box-size,
        height: box-size,
        align(box-align)[
          #rotate(box-rotate)[#counter(page).display("1")]
        ],
      ),
    )
  },
)

== The Anatomy of the Numinous

The exploration of human consciousness frequently encounters boundaries where
ordinary language begins to falter. It is within these liminal spaces that the
concept of the numinous finds its anchor, representing an experience that is
fundamentally distinct from the mundane structures of daily existence. When an
individual encounters the numinous, they are not merely confronting an
intellectual puzzle or a novel sensory stimulus; rather, they are intersecting
with a reality that feels entirely otherworldly, overwhelming, and deeply
saturated with an intrinsic, unnamable significance. This phenomenon bypasses
the rational faculties, striking directly at the emotional and existential core
of the observer.

#img(
  "../images/brains_out.jpg",
  overlay: gray.transparentize(70%),
)

Historically, the investigation into these profound states of awareness has
required a careful balance between psychological scrutiny and philosophical
openness. To categorize the experience without stripping away its vital essence
is the primary challenge of any such discourse. It demands an acknowledgment
that certain dimensions of perception operate outside the clean parameters of
empirical verification, existing instead as powerful subjective truths that
shape cultures, inspire artistic revolutions, and redefine individual life
trajectories across epochs.


=== Theoretical Foundations and Frameworks

To construct a coherent framework for analyzing these encounters, one must
dissect the constituent elements that define the underlying psychological state.
The primary layer of this experience is characterized by a profound sense of
dependence, a feeling where the self shrinks in the presence of an immense,
encompassing reality. This is not a destructive diminishment but a
transformative one, where the boundaries of the ego dissolve to reveal a broader
tapestry of existence. Researchers have noted that this initial shock to the
system serves as a necessary clearing of conventional thought patterns,
preparing the mind for deeper integration.

Following this initial dissolution, a secondary phase typically emerges, marked
by a fascinating yet terrifying attraction to the unknown. This duality ensures
that the observer is neither entirely repelled by the overwhelming scale of the
encounter nor completely absorbed by it to the point of losing all critical
faculties. Instead, a tense, dynamic equilibrium is maintained, allowing the
individual to remain conscious of their own position relative to the expansive
phenomenon unfolding before them. It is this precise tension that distinguishes
the truly numinous from simple aesthetic appreciation or ordinary emotional
resonance.

#img(
  "../images/kodak_vision3_200t_kodak_2393.jpg",
  overlay: gray.transparentize(70%),
)

==== Micro-analysis of Subjective Resonance

At the most granular level of experience, the cognitive faculties undergo a
temporary realignment. Standard temporal perception often warps, making moments
feel both instantaneous and eternal, a hallmark of deep psychological shifts.
The subject reports a heightened state of sensory clarity, where even the
silence between thoughts becomes heavy with potential meaning. This micro-level
shift suggests that the brain is processing information through channels that
are usually suppressed by the pragmatic demands of survival and routine.

Furthermore, the linguistic isolation that follows such an event cannot be
understated. Because the vocabulary of modern life is optimized for the exchange
of functional, material data, it lacks the necessary metaphors to convey a state
that is inherently non-rational. Individuals often resort to negative
descriptions, defining the experience by what it was not, or relying heavily on
paradoxical statements to hint at the complex truth they witnessed. This
communication barrier highlights the specialized nature of the phenomenon,
cementing its status as an isolated domain of human interiority.

== Manifestations Across Contexts

While the internal architecture of the experience remains remarkably consistent,
the outward forms it assumes are highly dependent on environmental and cultural
contexts. In some settings, the numinous is channeled through architectural
achievements, where soaring vaults, intentional geometry, and the deliberate
manipulation of light and shadow conspire to evoke a sense of vastness. In other
contexts, it is the raw, untamed expanse of the natural world that triggers the
shift, demonstrating that the human mind can find this profound resonance in
both deliberate human creation and chaotic ecological systems.

The societal implications of these shared experiences are vast, often serving as
the foundational myths or core values that bind communities together. When a
group collectively acknowledges a specific site, narrative, or symbol as a
repository of this non-rational power, they establish a sacred geography that
guides their ethical and social systems. This collective anchoring provides
stability, offering a shared reference point that transcends the shifting
political and economic currents of the era, though it also risks stagnation if
the living experience hardens into mere dogma.

=== Technological and Modern Intersections

In the contemporary era, the rapid acceleration of technological capability has
introduced a novel variable into this ancient dynamic. As systems grow more
complex, operating on scales of time and data that exceed human comprehension, a
secular variation of the numinous has begun to manifest within digital
environments. The vastness of interconnected networks and the emergence of
intricate computational models evoke a familiar blend of awe and apprehension,
suggesting that the drive to encounter the immense is not bound to traditional
frameworks but will adapt to whatever medium dominates the age.

This technological shift raises critical questions about the future of human
perception. If the non-rational mind can find a sense of overwhelming scale
within synthetic constructs, the line between organic reality and artificial
depth begins to blur. The challenge for modern observers is to discern whether
these digital echoes possess the same transformative potential as their
historical counterparts, or if they merely simulate the surface aesthetics of
awe without providing the profound existential grounding that characterized
traditional encounters.
