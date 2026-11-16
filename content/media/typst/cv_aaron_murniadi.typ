#set document(title: "Aaron P. Murniadi's CV")

#let section-block(title, content) = [
  #text(size: 1.25em, style: "italic", title)

  #block(inset: (left: 2em))[
    #content
  ]

  #v(0.5em)
]

#let work-heading(company: [], location: [], time: []) = [
  #grid(
    columns: (3fr, 2fr),
    align: (left, right),
    [
      #text(weight: "bold")[#smallcaps[#company]]
      #if location != [] [| #location]
    ],
    time,
  )
]

#let role-heading(role: [], time: []) = [
  #v(0.8em)
  #grid(
    columns: (1fr, auto),
    align: (left + bottom, right + bottom),
    text(style: "italic", weight: "bold", role + " (" + time + ")"), v(0.2em),
  )
]

#let edu-heading(department: [], location: [], role: [], time: []) = [
  #grid(
    columns: (3fr, 1fr),
    align: (left, right),
    row-gutter: 1em,
    [#strong(department) \ #emph(role)], [#location \ #time],
  )
]

#let publication-heading(detail: [], year: []) = [
  #grid(
    columns: (5fr, 1fr),
    align: (left, right),
    detail, year,
  )
]

#let header(name: [], contact: []) = [
  #grid(
    columns: (3fr, 2fr),
    align: (left + bottom, right + bottom),
    text(size: 2em, style: "italic")[#name], [#contact],
  )
  #line(length: 100%, stroke: 1pt)
]

#let point = text(size: 0.95em, weight: "bold")[★]

// Hanging-indent bullet list. Each item is its own line, so long entries wrap
// under the text and not under the star.
#let bullets(..items) = {
  set par(justify: true, leading: 0.55em)
  for item in items.pos() {
    block(
      above: 0.5em,
      below: 0pt,
      grid(
        columns: (1.1em, 1fr),
        align: (left + top, left + top),
        point, item,
      ),
    )
  }
}

// =============================================
// =============================================

#let today = datetime.today().display("[month repr:long] [day], [year]")

#set par(
  justify: true,
  justification-limits: (
    tracking: (max: 0.02em, min: -0.015em),
  ),
)

#set text(
  hyphenate: true,
  lang: "en",
  number-width: "proportional",
  costs: (hyphenation: 30%),
)

#set page(paper: "a4", margin: 0.95in, footer: [#text(
    style: "italic",
    size: 1em,
  )[Last updated: #today #h(1fr)] #context [\[#counter(page).display() / #(
      counter(page).final().first()
    )]\]])

#header(
  name: [Aaron P. Murniadi],
  contact: [
    #link("mailto:aaronmurniadi@gmail.com")[_aaronmurniadi\@gmail.com_] \
    (+62) 881 535 9522, Jakarta, Indonesia \
    #link("https://www.linkedin.com/in/aaronmurniadi/")[LinkedIn] | #link(
      "https://github.com/aaronmurniadi",
    )[GitHub] | #link("https://aaron.beago-cirius.ts.net/")[My Blog]
  ],
)

Software engineer with more than four years of experience across the full stack
of a production system: operating system images, a reverse engineered device
driver, Go backend services, web front ends, and test automation. Replaced a
proprietary vendor runtime with an in-house driver, rewrote the capture path in
Go to halve the processing time per unit, and built the tooling that provisions
and monitors a 220 machine fleet across four countries. Previous background in
test automation and security research.

#section-block(
  [Work Experience],
  [#work-heading(
      company: [#link("https://www.ninjavan.co/")[Ninja Van]],
      location: [Jakarta, Indonesia],
      time: [December 2021--_Present_],
    )

    #text(size: 0.9em, style: "italic")[Promoted three times in four years, from
      Software Tester, Automation Engineer, to Software Engineer 2.]

    #role-heading(role: [Software Engineer 2, Developer], time: [July
      2026--_Present_])

    #bullets(
      [Reverse engineered the proprietary camera SDK and wire protocol and wrote
        a replacement machine vision driver. Frame rate rose from 2 to 8 frames
        per second, and the vendor runtime cost dropped to zero.],
      [Cut scan time per parcel from 3.0 s to 1.5 s, with a 2.0 s recovery pass
        on parcels that the camera does not detect, by rewriting the capture
        path in Go.],
      [Own the DWS software stack for 220 production machines in Indonesia,
        Malaysia, Singapore, and the Philippines, with a combined peak scan rate
        of 15,000 parcels per hour across shifts.],
      [Built the Debian based operating system image for DWS hardware with the
        Debian Live ISO builder: 2.5 GB image, under 30 s to ready state, 25 min
        unattended install.],
      [Built the capture application with Wails (Go and SvelteKit), which runs
        at up to 10,000 parcels per machine per day, and the provisioning
        dashboard that 15 staff in four countries use to roll out the fleet.],
      [Assigned 10 million latitude and longitude waypoints to hundreds of
        thousands of waypoint zones with PostGIS for parcel routing.],
      [Resolve about 5 production incidents per month within one hour on a 24/7
        on-call rotation of 6 engineers.],
    )

    #role-heading(
      role: [Software Engineer 2, QA],
      time: [September 2024--July 2025],
    )

    #bullets(
      [Planned and led the migration of the Karate regression projects from
        Maven to Gradle, covering 5 teams and about 1,500 scenarios.],
      [Built a multi-scale template matching tool with #link(
          "https://opencv.org/",
        )[OpenCV] that finds UI elements from reference screenshots, which
        removed manual element lookup from the front-end suites.],
    )

    #role-heading(
      role: [Software Engineer 1, QA],
      time: [January 2023--August 2024],
    )

    #bullets(
      [Designed and maintained the Java regression suite for 5 teams: API tests
        with #link("https://rest-assured.io/")[REST Assured] and #link(
          "https://square.github.io/okhttp/",
        )[OkHttp], web tests with #link(
          "https://www.selenium.dev/documentation/webdriver/",
        )[Selenium WebDriver], and Android tests with #link(
          "https://github.com/appium/appium-flutter-driver",
        )[Appium Flutter Driver], all on #link(
          "https://github.com/cucumber/cucumber-jvm",
        )[Cucumber-JVM].],
      [Led the design and implementation of an API to Scenario (ATS) mapping
        service in Go that traces each automated scenario to the backend APIs it
        covers.],
    )

    #role-heading(
      role: [Software Tester (Internship)],
      time: [December 2021--December 2022],
    )

    #bullets(
      [Ran manual test passes on new and existing features against their
        requirements.],
      [Built a Python #link("https://github.com/behave/behave")[Behave] harness
        that gave manual testers semi-automation, used by 5 colleagues.],
    )
  ],
)

#block(breakable: false, section-block([Technical Skills], [
  #set text(hyphenate: false)
  #table(
    columns: (auto, auto),
    inset: 5pt,
    stroke: 0pt,
    fill: (none, none),
    align: (right, left),
    [*Programming Languages*], [Java, Go, Python, TypeScript, JavaScrip, C++],

    [*Application Development*],
    [React, SvelteKit, Wails, Yarn, #link(
        "https://swagger.io/specification/",
      )[OpenAPI]],

    [*Testing*],
    [
      #link("https://www.selenium.dev/")[Selenium], #link(
        "https://appium.io/",
      )[Appium], #link("https://github.com/karatelabs/karate")[Karate DSL],
      #link("https://github.com/cucumber/cucumber-jvm")[Cucumber-JVM], #link(
        "https://rest-assured.io/",
      )[REST Assured], #link("https://github.com/behave/behave")[Behave], #link(
        "https://square.github.io/okhttp/",
      )[OkHttp],
      #link("https://opencv.org/")[OpenCV]
    ],

    [*Data and Infrastructure*],
    [
      #link("https://postgis.net/")[PostGIS], #link(
        "https://www.mysql.com/",
      )[MySQL], #link("https://redis.io/")[Redis], #link(
        "https://docs.docker.com/",
      )[Docker], #link(
        "https://live-team.pages.debian.net/live-manual/html/live-manual/index.en.html",
      )[Debian Live ISO Builder],#link(
        "https://www.atlassian.com/software/bamboo",
      )[Bamboo], #link("https://gradle.org/")[Gradle], #link(
        "https://maven.apache.org/",
      )[Maven],
      #link("https://git-scm.com/")[Git]
    ],

    [*Operating Systems*],
    [
      #link("https://www.microsoft.com/windows/")[Windows], #link(
        "https://www.apple.com/macos/",
      )[macOS], #link("https://www.debian.org/")[Debian], #link(
        "https://www.fedora.org/",
      )[Fedora], #link("https://ubuntu.com/")[Ubuntu], #link(
        "https://archlinux.org/",
      )[Arch Linux]
    ],

    [*Other*],
    [Reverse engineering, machine vision, agentic code generation workflows,
      technical writing],
  )
]))

#section-block(
  [Security Research],
  [
    Vulnerability researcher, 2022--2023. Both findings were disclosed privately
    to the vendor and remain under NDA.

    #v(0.4em)
    #bullets(
      [Directory traversal in a journal publication platform, which exposed full
        text PDFs behind a paid access wall.],
      [Query enumeration and Insecure Direct Object Reference (IDOR) in a
        national education institution system, which exposed personally
        identifiable information (PII).],
    )
  ],
)

#section-block(
  [Projects],
  [
    #bullets(
      [*#link("https://github.com/aaronmurniadi/gogige")[gogige]*: a pure Go
        GigE Vision client (GVCP control channel, GVSP streaming with packet
        resend, GenICam GenApi node map) that drives industrial machine vision
        cameras without a vendor SDK, driver, or GenTL producer. BSD-3-Clause,
        v1.10.0.],
      [Built the implementation from the standard documents: ran the GigE Vision
        and GenICam reference PDFs through Mistral OCR into structured markdown,
        then developed an agent and subagent workflow that turns the parsed
        reference into Go code.],
    )
  ],
)

#section-block(
  [Education],
  [#edu-heading(
      department: [SMAK 1 Penabur Cirebon],
      location: [Cirebon, West Java],
      role: [Natural Sciences Concentration],
      time: [2014--2017],
    )
    Finalist, ANFORCOM (Annual Informatics Competition) 2017, #link(
      "https://www.undip.ac.id/",
    )[Diponegoro University], Semarang.
  ],
)

#section-block(
  [Publications],
  [
    #publication-heading(
      detail: [Murniadi, Aaron Purnomo (2022) "Waktu Intuitif (durasi) sebagai
        Dasar Adanya Kehendak Bebas Manusia Menurut Henri Bergson." #link(
          "https://karya.brin.go.id/id/eprint/25247",
        )[_Driyarkara Jurnal Filsafat_], 42 (2): 5. pp. 81--91. ISSN 2809-9516],
      year: [2022],
    )
    #publication-heading(
      detail: ["Ranting Mencari Jalan ke Akar", in Mardiatmadja, SJ, B.S. and
        Whisnu Bintoro, CICM, Dhaniel (2020) #link(
          "http://repo.driyarkara.ac.id/321/1/EKLESIOLOGI.PDF",
        )[_Eklesiologi Langkah demi Langkah_], PT Kanisius, Yogyakarta. ISBN
        978-979-21-6665-1],
      year: [2020],
    )
  ],
)
