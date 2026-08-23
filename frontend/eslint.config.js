import js from '@eslint/js'
import vue from 'eslint-plugin-vue'
import tseslint from '@vue/eslint-config-typescript'

export default [
  { ignores: ['dist/**', 'coverage/**'] },
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  ...tseslint(),
  {
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/max-attributes-per-line': 'off',
      'vue/html-closing-bracket-newline': 'off',
      'vue/html-indent': 'off',
      'vue/multiline-html-element-content-newline': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/html-self-closing': 'off',
    },
  },
]
