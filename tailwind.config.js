module.exports = {
  content: [
    './app/v2/views/**/*.templ',
    './app/v2/views/**/*_templ.go',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Lufga", "Syne", "Nunito", "sans-serif", "ui-sans-serif", "system-ui"],
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/aspect-ratio"),
    require("daisyui"),
  ],
  daisyui: {
    themes: [
      {
        memnix: {
          primary: "#E9AF98",
          secondary: "#405CA0",
          accent: "#AD6E9E",
          neutral: "#40404A",
          "base-200": "#f1f1f1",
          "base-100": "#FCFCFC",
          info: "#3A73D4",
          success: "#72E9C1",
          warning: "#F6D73C",
          error: "#E83B55",
          "primary-content": "#4d1600",
          "success-content": "#102742",
          "error-content": "#000000",
        },
      },
      {
        dark: {
          primary: "#E9AF98",
          secondary: "#405CA0",
          accent: "#AD6E9E",
          neutral: "#40404A",
          "base-100": "#2A303C",
          "base-200": "#242933",
          info: "#3A73D4",
          success: "#72E9C1",
          warning: "#F6D73C",
          error: "#E83B55",
          "primary-content": "#4d1600",
          "success-content": "#102742",
          "error-content": "#000000",
        },
      },
    ],
  },
};
