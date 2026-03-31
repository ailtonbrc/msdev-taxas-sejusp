import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/taxas/',
  
  // 1. Configurações para quando você roda "npm run dev" (No seu computador)
  server: {
    port: 3000, 
    proxy: {
      '/api': {
        target: 'http://localhost:4001', // Backend agora está na porta 4001
        changeOrigin: true,
        secure: false,
      },
    },
  },

  // 2. Configurações para quando você roda "npm run build" (Para o servidor)
  build: {
    outDir: 'dist',
    emptyOutDir: true, // Limpa a pasta dist antes de gerar novos arquivos
  }
});