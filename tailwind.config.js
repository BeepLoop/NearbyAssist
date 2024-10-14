/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.{templ, go}"],
    theme: {
        extend: {
            colors: {
                primary: {
                    DEFAULT: '#4ECB71',
                    pale: '#DCFFB7',
                    greener: '#163020'
                },
            }
        },
    },
    plugins: [],
}

