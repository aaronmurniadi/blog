const pluginRss = require("@11ty/eleventy-plugin-rss");
const pluginSyntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const pluginBundle = require("@11ty/eleventy-plugin-bundle");
const pluginSeo = require("eleventy-plugin-seo");
const markdownIt = require("markdown-it");
const markdownItAnchor = require("markdown-it-anchor");
const markdownItAttrs = require("markdown-it-attrs");
const markdownItFootnote = require("markdown-it-footnote");
const eleventyPluginFilesMinifier = require("@sherby/eleventy-plugin-files-minifier");
const path = require("node:path");
const fs = require("node:fs/promises");
const CleanCSS = require("clean-css");
const htmlMinifier = require("html-minifier");

/** Match @sherby/eleventy-plugin-files-minifier HTML options (passthrough HTML skips transforms). */
const HTML_MINIFY_OPTIONS = {
  collapseBooleanAttributes: true,
  collapseWhitespace: true,
  decodeEntities: true,
  html5: true,
  minifyCSS: true,
  minifyJS: true,
  removeComments: true,
  removeEmptyAttributes: true,
  removeEmptyElements: false,
  sortAttributes: true,
  sortClassName: true,
  useShortDoctype: true
};

async function walkHtmlFiles(dir) {
  const files = [];
  const entries = await fs.readdir(dir, { withFileTypes: true });
  for (const ent of entries) {
    const full = path.join(dir, ent.name);
    if (ent.isDirectory()) {
      files.push(...(await walkHtmlFiles(full)));
    } else if (ent.isFile() && ent.name.endsWith(".html")) {
      files.push(full);
    }
  }
  return files;
}

module.exports = function (eleventyConfig) {
  // Plugins
  eleventyConfig.addPlugin(eleventyPluginFilesMinifier);
  eleventyConfig.addPlugin(pluginRss);
  eleventyConfig.addPlugin(pluginSyntaxHighlight, {
    preAttributes: {
      tabindex: -1
    },
    alwaysWrapLineHighlights: false,
    trim: true
  });
  eleventyConfig.addPlugin(pluginBundle);
  eleventyConfig.addPlugin(pluginSeo);

  // Custom shortcode for rendering post lists
  eleventyConfig.addShortcode("postList", function (collection) {
    let items = collection.map(item => {
      const date = item.date ? `<span>${new Date(item.date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })} — </span>` : '';
      return `<li>${date}<a href="${item.url}">${item.data.title || item.title}</a></li>`;
    }).join('');
    return `<ul>${items}</ul>`;
  });

  // Markdown configuration
  let markdownLib = markdownIt({
    html: true,
    breaks: false,
    linkify: true,
    typographer: true,
    quotes: ['\u201c', '\u201d', '\u2018', '\u2019']
  })
    .use(markdownItAnchor, {
      permalink: markdownItAnchor.permalink.linkInsideHeader({
        placement: "after",
        symbol: "#",
        space: true,
        renderAttrs: () => ({
          "aria-label": "Permalink to this section",
        }),
      }),
    })
    .use(markdownItAttrs)
    .use(markdownItFootnote);

  eleventyConfig.setLibrary("md", markdownLib);

  // Markdown library without anchors for simple layouts
  let markdownLibNoAnchors = markdownIt({
    html: true,
    breaks: true,
    linkify: true,
    typographer: true,
    quotes: ['\u201c', '\u201d', '\u2018', '\u2019']
  })
    .use(markdownItAttrs)
    .use(markdownItFootnote);

  eleventyConfig.setLibrary("md", markdownLib);
  eleventyConfig.addPairedShortcode("markdownNoAnchors", function (content) {
    return markdownLibNoAnchors.render(content);
  });

  // Collections
  eleventyConfig.addCollection("posts", function (collectionApi) {
    return collectionApi.getFilteredByGlob("_posts/**/*.md").sort(function (a, b) {
      return b.date - a.date;
    });
  });

  eleventyConfig.addCollection("articles", function (collectionApi) {
    return collectionApi.getFilteredByGlob("_articles/**/*.md").sort(function (a, b) {
      return b.date - a.date;
    });
  });

  eleventyConfig.addCollection("recaps", function (collectionApi) {
    return collectionApi.getFilteredByGlob("_recaps/**/*.md").sort(function (a, b) {
      return b.date - a.date;
    });
  });

  eleventyConfig.addCollection("summaries", function (collectionApi) {
    return collectionApi.getFilteredByGlob("_summaries/**/*.md").sort(function (a, b) {
      return b.date - a.date;
    });
  });

  eleventyConfig.addCollection("tools", function (collectionApi) {
    return collectionApi.getFilteredByGlob("_tools/**/*.md").sort(function (a, b) {
      return b.date - a.date;
    });
  });

  // Sitemap collection
  eleventyConfig.addCollection("sitemap", function (collectionApi) {
    const allPages = [
      // Add main pages
      { url: "/", date: new Date(), changefreq: "weekly", priority: "1.0" },
      { url: "/contact/", date: new Date(), changefreq: "monthly", priority: "0.8" },
      { url: "/posts/", date: new Date(), changefreq: "weekly", priority: "0.9" },
      { url: "/articles/", date: new Date(), changefreq: "weekly", priority: "0.9" },
      { url: "/recaps/", date: new Date(), changefreq: "weekly", priority: "0.9" },
      { url: "/summaries/", date: new Date(), changefreq: "weekly", priority: "0.9" },
      { url: "/toolbox/", date: new Date(), changefreq: "monthly", priority: "0.7" },
      { url: "/photography/", date: new Date(), changefreq: "monthly", priority: "0.7" },
      { url: "/typesettings/", date: new Date(), changefreq: "monthly", priority: "0.7" }
    ];

    // Add all collection items
    const collections = ["posts", "articles", "recaps", "summaries", "tools"];
    collections.forEach(collectionName => {
      const collection = collectionApi.getFilteredByGlob(`_${collectionName}/**/*.md`);
      collection.forEach(item => {
        allPages.push({
          url: item.url,
          date: item.date,
          changefreq: "monthly",
          priority: "0.8"
        });
      });
    });

    return allPages.sort((a, b) => a.url.localeCompare(b.url));
  });

  // Copy assets
  eleventyConfig.addPassthroughCopy("assets");
  eleventyConfig.addPassthroughCopy("media");
  eleventyConfig.addPassthroughCopy("avatar.jpg");
  eleventyConfig.addPassthroughCopy("robots.txt");
  eleventyConfig.addPassthroughCopy("_includes/**/*.(css|js|jpg|jpeg|png|gif|svg|ico|pdf)");
  eleventyConfig.addPassthroughCopy("_tools");

  // Passthrough assets skip template transforms; minify CSS + all HTML after write.
  eleventyConfig.on("eleventy.after", async ({ directories }) => {
    const outDir = directories.output;

    const cssPath = path.join(outDir, "assets", "css", "style.css");
    try {
      const source = await fs.readFile(cssPath, "utf8");
      const { styles, errors } = new CleanCSS({ level: 2 }).minify(source);
      if (errors?.length) {
        console.warn("[clean-css] style.css:", errors.join("; "));
      }
      await fs.writeFile(cssPath, styles, "utf8");
    } catch (err) {
      if (err.code !== "ENOENT") throw err;
    }

    const htmlPaths = await walkHtmlFiles(outDir);
    for (const htmlPath of htmlPaths) {
      try {
        const raw = await fs.readFile(htmlPath, "utf8");
        const min = htmlMinifier.minify(raw, HTML_MINIFY_OPTIONS);
        await fs.writeFile(htmlPath, min, "utf8");
      } catch (err) {
        console.warn("[html-minifier]", path.relative(outDir, htmlPath), err.message);
      }
    }
  });

  // Date filters
  eleventyConfig.addFilter("dateReadable", dateObj => {
    return new Date(dateObj).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  });

  eleventyConfig.addFilter("dateToRfc3339", dateObj => {
    return new Date(dateObj).toISOString();
  });

  eleventyConfig.addFilter("date", (dateObj, format) => {
    const date = new Date(dateObj);
    if (format === "%Y") {
      return date.getFullYear();
    }
    return date.toISOString();
  });

  // Global data
  eleventyConfig.addGlobalData("site", {
    title: "Beago Cirius",
    email: "aaronmurniadi@gmail.com",
    description: "A Software Engineer with a passion for technology, coding, and open-source software.",
    baseurl: "",
    url: "https://aaron.beago-cirius.ts.net",
    lang: "en-US",
    github_username: "aaronmurniadi",
  });

  return {
    dir: {
      input: ".",
      includes: "_includes",
      data: "_data",
      output: "_site",
      layouts: "_layouts"
    },
    templateFormats: ["md", "njk", "html"],
    htmlTemplateEngine: "njk",
    markdownTemplateEngine: "njk",
    dataTemplateEngine: "njk",
    passthroughFileCopy: true,
    cleanUrls: true
  };
};
