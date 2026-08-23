# Frontend agent guide

- Build with Vue 3 Composition API, TypeScript, Vite, and `<script setup>`.
- Keep Pinia limited to authenticated session and favorite clubs.
- Prefer semantic HTML, keyboard-visible controls, clear focus states, and reduced-motion-safe transitions.
- Reuse design tokens in `src/style.css`; do not hardcode ad-hoc theme colors in components.
- Use Reka UI primitives for complex interactive widgets and Lucide Vue for interface icons.
- Cover user-visible behavior with Vitest; keep API access behind `src/lib/api.ts`.
- Validate layouts at 375 px and a desktop width, in both light and dark themes.
