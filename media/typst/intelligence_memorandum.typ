#let meta = (
  doc_type: "INTELLIGENCE MEMORANDUM",
  doc_date: "14 April 1966",
  oci_no: "1170/66",
  copy_no: "236",
  title: "INDONESIA'S \"1945 CONSTITUTION\"",
  warning: "",
  body: [
    1\. When the present Indonesian constitution was written in 1945, it was designed more as a symbol of the struggle for national independence than as a document defining a system of government for a modern state. It was superseded in 1949, and a relatively unsuccessful experiment in parliamentary democracy ensued. In 1959, however, Sukarno resurrected the 1945 constitution to legitimatize his authoritarian regime of "guided democracy." Under this extremely vague document, great power is vested in the executive, and there is almost no guarantee for a popular voice in the government.

    2\. Military and civilian leaders of the new government are now criticizing Sukarno's abuse of presidential power and calling for a "return" to constitutional principles. This agitation, however, appears largely a tactic to further the denigration of Sukarno. The army seems likely to use the 1945 constitution to justify a relatively authoritarian regime much as Sukarno did.

    3\. Under the 1945 constitution, the president is chosen by and is responsible only to the People's Consultative Congress (MPR), a body composed of representatives of political and functional groups and charged with "deciding the outlines of national policy." The constitution merely states that the membership of the MPR, which must meet at least once every five years, is to be chosen in a manner "pre-scribed by law." No such law has materialized and to date only an "interim" MPR appointed by Sukarno has met.

    4\. Legislative power is shared by the president and the Council of Peoples Representatives (DPR) or parliament, a body whose members hold concurrent membership in the MPR. No method is prescribed for selecting these members, however, and they, too, have been hand-picked by Sukarno. Both the president and the DPR can initiate legislation but the DPR cannot override a presidential veto, and "at critical times" its legislation can be replaced by presidential decree. This one-sided relationship has been reinforced by a uniquely Indonesian order of procedure inaugurated by Sukarno and not stipulated by the constitution. Under this procedure, unless the DPR can achieve a unanimous vote, it must refer the matter to the president for final decision.

    5\. The constitution provides for a vice-president who "assists" the president and fills the presidency upon the incumbent's death or disability. Since the reinvocation of the constitution in 1959, however, the vice-presidency has been vacant. In 1963 Sukarno's hand-picked MPR appointed him president for life. Sukarno accepted the appointment with the understanding that it would be "reviewed" by the first popularly elected MPR.

    6\. Military and civilian leaders now plan to convene a session of the MPR in mid-May which they will no doubt use to contrast the new government's responsibility to the people with that of Sukarno's regime. Among other things this congress may fill the vice-presidency and might even revoke Sukarno's lifetime mandate, although this appears unlikely. Any significant "democratization" of the 1945 Constitution, however, seems unlikely. The army, long disenchanted with Indonesian political parties, has always approved Sukarno's concept of "guided democracy" and strongly supported his junking of the parliamentary system in 1959. As a result the army is likely to work through Sukarno's interim MPR rather than hold national elections in the near future. With pro-Communist and Communist MPR members arrested or "liquidated," the body can easily be manipulated.

    7\. In the political vacuum left by the eclipse of Sukarno, the army, in fact, has little choice in the immediate future but to maintain its authoritarian control. An almost inevitable point of difference between the military and its civilian allies will be the eventual extent of constitutional reform.
  ],
)


//  Page geometry 

#set page(
  paper: "us-gov-letter",
  margin: (top: 1.5in, bottom: 1in, left: 1.55in, right: 1.55in),
)

//  Typography 

// This is the root font size for the memorandum body.
#set text(
  size: 11pt,
  lang: "en",
  hyphenate: true,
  number-type: "lining",
)

#set par(
  first-line-indent: (amount: 3em, all: true),
  leading: 0.5em,
  spacing: 1.25em,
  justify: false,
)

#show ". ": ".  "

// COVER
#[

  #place(top + left, dx: -1cm, dy: -1cm)[
    #image("inteeligence_memorandum_logo.png", width: 3.5cm)
    #v(1fr)
  ]

  #place(
    top + left,
    dy: -60pt,
    dx: -66pt,
    line(length: 200%, stroke: 16pt),
  )
  #place(
    top + left,
    dy: -68pt,
    dx: -60pt,
    line(length: 200%, angle: 90deg, stroke: 16pt),
  )

  #align(right)[
    #block(width: 9em)[
      #table(
        columns: (2fr, 1fr),
        align: (left, center),
        stroke: none,
        inset: 3pt,
        [#meta.doc_date], [],
        [], [],
        [], [],
        [OCI No.], [#meta.oci_no],
        [COPY  No.], [#text(size: 18pt)[#meta.copy_no]],
      )]
  ]

  #v(1fr)

  #align(center)[
    #text(tracking: 1pt, weight: "bold", size: 16pt)[#meta.doc_type]
    #v(1fr)
    #text(tracking: 1pt, size: 12pt)[#meta.title]
    #v(2fr)
    #text(tracking: 1pt)[DIRECTORATE OF INTELLIGENCE]\ \
    Office of Current Intelligence
  ]
  #v(1fr)
]

#pagebreak()

// WARNING
#[
  #set par(first-line-indent: 0em)
  #align(center)[#block(width: 23em)[
    #v(2fr)
    WARNING
    #align(left)[
      This document holds data relevant to the national defense of the Federation of Valoria as defined by Title 18, Sections 793 and 794 of the Valorian Legal Code. Legal statutes forbid sharing, disclosing, or allowing an unauthorized individual to receive its contents. Additionally, making copies of this form is strictly prohibited.
      #v(1fr)]
  ]]
]#pagebreak()

// BODY
#[#set text(font: "PT Mono")
  #align(right)[OCI NO. #meta.oci_no]\

  #align(center)[
    MYSTICAL INVESTIGATION AGENCY \
    Office of Current Intelligence \
    #meta.doc_date
  ]

  \
  #meta.doc_type \
  \

  #align(center)[
    #underline(
      stroke: 1pt, // Thicker red line
      offset: 2pt, // Distance from the text
      [#meta.title],
    )
  ]

  #meta.body
]

