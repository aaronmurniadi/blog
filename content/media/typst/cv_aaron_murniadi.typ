#set document(title: "Aaron P. Murniadi's CV", author: "Aaron P. Murniadi")

// ---------- Page and text ----------
#set page(
  paper: "a4",
  margin: (top: 0.7in, x: 0.8in, bottom: 0.7in),
  footer: context align(center, text(size: 0.85em)[
    #counter(page).display() / #counter(page).final().first()
  ]),
)

#set text(size: 10.5pt, lang: "en", hyphenate: false)
#set par(justify: false, leading: 0.72em, spacing: 0.9em)
#set list(marker: [•], indent: 0.3em, body-indent: 0.6em, spacing: 0.65em)

// ---------- Helpers ----------
#let section-block(title, content) = [
  #block(above: 1.5em, below: 0.8em, sticky: true)[
    #text(size: 1.05em, weight: "bold", tracking: 0.04em)[#upper(title)]
    #v(-0.3em)
    #line(length: 100%, stroke: 0.5pt)
  ]
  #content
]

#let work-heading(company: [], location: [], time: []) = grid(
  columns: (1fr, auto),
  align: (left, right),
  [#text(weight: "bold")[#company] | #location | #time], [],
)

#let role-heading(role: [], time: []) = block(
  above: 1.2em,
  below: 0.7em,
  sticky: true,
)[
  #grid(
    columns: (1fr, auto),
    align: (left, right),
    text(weight: "bold", style: "italic")[#role (#time)],
    text(style: "italic")[],
  )
]

#let entry-heading(main: [], side: []) = grid(
  columns: (1fr, auto),
  align: (left, right),
  main, side,
)

#let bullets(..items) = list(..items.pos())

// ---------- Header ----------
#align(center)[
  #text(size: 1.9em, weight: "bold")[Aaron P. Murniadi] \
  #v(0.2em)
  Jakarta, Indonesia | (+62) 881 535 9522 | #link(
    "mailto:aaronmurniadi@gmail.com",
  )[aaronmurniadi\@gmail.com] \
  #link(
    "https://www.linkedin.com/in/aaronmurniadi/",
  )[linkedin.com/in/aaronmurniadi] |
  #link("https://github.com/aaronmurniadi")[github.com/aaronmurniadi]
]

// ---------- Summary ----------
#section-block([Summary], [
  Software engineer with 5+ years of experience across the full stack of a
  production system: OS images, a reverse engineered device driver, Go and Java
  backend services, web front ends, and monitoring. Replaced a proprietary
  vendor runtime with an in-house driver, halved per-parcel scan time by
  rewriting the capture path in Go, and built the tooling that provisions,
  monitors, and hardens a 220 machine fleet in four countries. Previous
  background in test automation and security research.
])

// ---------- Experience ----------
#section-block([Work Experience], [
  #work-heading(
    company: [Ninja Van],
    location: [Jakarta, Indonesia],
    time: [Dec 2021 – Dec 2026],
  )
  #text(size: 0.92em, style: "italic")[Promoted three times in five years:
    Software Tester, Automation Engineer, Software Engineer 2.]

  #role-heading(role: [Software Engineer 2, Developer], time: [Aug 2025 – Dec
    2026])
  #bullets(
    [Own the DWS software stack for 220 production machines in Indonesia,
      Malaysia, Singapore, and the Philippines, at a combined peak of 15,000
      parcels per hour.],
    [Reverse engineered the proprietary camera SDK and GigE Vision protocol
      (Wireshark) and wrote a replacement machine vision driver in Go
      (libaravis, CGO): frame rate rose from 2 to 8 FPS and the vendor
      maintenance cost dropped to zero.],
    [Rewrote the capture path in Go, cutting scan time per parcel from 3.0 s to
      1.5 s (2.0 s recovery pass for parcels the camera does not detect).],
    [Built the Debian based OS image with the Debian Live ISO builder (2.5 GB,
      under 30 s to ready state, 25 min unattended install). Hardened it for
      InfoSec review with Wazuh, RustDesk behind an nftables guard, and
      device-specific host names.],
    [Built the capture app with Wails (Go, SvelteKit), running up to 10,000
      parcels per machine per day, and the provisioning dashboard used by 15
      staff in four countries.],
    [Built fleet observability: two Grafana dashboards over Quickwit logs and
      Prometheus metrics (one with 30 panels), an ingest endpoint in the
      hub-users Go service, and a log uploader in the capture app.],
    [Resolve about 5 production incidents per month within one hour on a 24/7
      on-call rotation of 6 engineers. Mentored a new developer on the DWSv2
      architecture and found a silent crash in the Mile 0 order sync script.],
  )

  #role-heading(role: [Software Engineer 2, QA], time: [Jul 2024 – Aug 2025])
  #bullets(
    [Release manager for Driver App, coordinating three product teams with no
      rollbacks or blockers.],
    [Built the NVQA Tools IntelliJ plugin: Bamboo build monitoring with rerun, a
      test data encryption helper, and an AI chat over QA scenarios and
      Confluence pages.],
    [Led the migration of Karate regression projects from Maven to Gradle across
      5 teams and about 1,500 scenarios.],
    [Built a multi-scale template matching tool with OpenCV that finds UI
      elements from reference screenshots, removing manual element lookup from
      front-end suites.],
  )

  #block(breakable: false)[
    #role-heading(role: [Software Engineer 1, QA], time: [Jan 2023 – Jul 2024])
    #bullets(
      [Released common-karate, the shared automation library adopted by 5 teams,
        and supported their migrations.],
      [Maintained the Java regression suite for 5 teams (REST Assured, OkHttp,
        Selenium WebDriver, Appium Flutter Driver, Cucumber-JVM).],
      [Led design of an API to Scenario (ATS) mapping service in Go that traces
        each automated scenario to the backend APIs it covers.],
    )
  ]

  #block(breakable: false)[
    #role-heading(role: [Software Tester Intern], time: [Dec 2021 – Dec 2022])
    #bullets(
      [Ran manual test passes against requirements and built a Python Behave
        harness for semi-automation, used by 5 colleagues.],
    )
  ]
])

// ---------- Skills ----------
#section-block([Technical Skills], [
  #set par(spacing: 0.8em)
  *Programming Languages:* Java, Go, Python, C++ \
  *Development Frameworks:* React, SvelteKit, Wails, Yarn \
  *Testing:* Selenium, Appium, Karate DSL, Cucumber-JVM, REST Assured, Behave,
  OkHttp, OpenCV, OpenAPI \
  *Data and Infrastructure:* PostGIS, MySQL, Redis, Elasticsearch, Docker,
  Debian Live ISO Builder, Bamboo, Gradle, Maven, Git, Grafana, Prometheus,
  Quickwit, Wazuh \
  *Operating Systems:* Debian, Ubuntu, Fedora, Arch Linux, Windows, macOS \
  *Other:* Reverse engineering, machine vision, agentic code generation
  workflows, technical writing
])

// ---------- Projects ----------
#section-block([Projects], [
  #bullets(
    [*#link("https://github.com/aaronmurniadi/gogige")[gogige]*: pure Go GigE
      Vision client (GVCP control, GVSP streaming with packet resend, GenICam
      node map) that drives industrial cameras without a vendor SDK or GenTL
      producer. BSD-3-Clause, v1.10.0. Built from the standard by OCR-ing the
      reference PDFs into markdown and using an agent workflow to generate Go
      code.],
    [*#link("https://github.com/aaronmurniadi/ofxhost")[OfxHost]*: A minimal
      still-image OpenFX plugin host. Open a RAW (or PNG/JPEG/TIFF/EXR) photo,
      run it through one or more OFX filter plugins, preview the result, and
      export to PNG, JPEG, TIFF (8/16-bit integer or 32-bit float).],
  )
])

// ---------- Security ----------
#section-block([Security Research], [
  Vulnerability researcher, 2022 – 2023: found a directory traversal in a
  journal publication platform and an IDOR with query enumeration in a national
  education system. Both were disclosed privately to the vendor and remain under
  NDA.
])

// ---------- Education ----------
#section-block([Education], [
  #entry-heading(
    main: [*SMAK 1 Penabur Cirebon*, Natural Sciences],
    side: [Cirebon, West Java | 2014 – 2017],
  )
  Finalist, ANFORCOM 2017, Diponegoro University, Semarang.
])

// ---------- Publications ----------
#section-block([Publications], [
  #entry-heading(
    main: [Murniadi, A. P. "Waktu Intuitif (durasi) sebagai Dasar Adanya
      Kehendak Bebas Manusia Menurut Henri Bergson." #link(
        "https://karya.brin.go.id/id/eprint/25247",
      )[_Driyarkara Jurnal Filsafat_], 42 (2), pp. 81–91. ISSN 2809-9516.],
    side: [2022],
  )
  #v(0.6em)
  #entry-heading(
    main: ["Ranting Mencari Jalan ke Akar," in _Eklesiologi Langkah demi
      Langkah_, PT Kanisius, Yogyakarta. ISBN 978-979-21-6665-1.],
    side: [2020],
  )
])
