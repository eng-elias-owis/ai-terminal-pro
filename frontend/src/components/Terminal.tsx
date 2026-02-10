import { useEffect, useRef } from 'react'
import { Terminal as XTerm } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

interface TerminalProps {
  onOpenAI: () => void
}

export const Terminal: React.FC<TerminalProps> = ({ onOpenAI }) => {
  const terminalRef = useRef<HTMLDivElement>(null)
  const xtermRef = useRef<XTerm | null>(null)
  const fitAddonRef = useRef<FitAddon | null>(null)

  useEffect(() => {
    if (!terminalRef.current) return

    // Initialize xterm.js
    const term = new XTerm({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'JetBrains Mono, monospace',
      theme: {
        background: '#1a1a1a',
        foreground: '#e0e0e0',
        cursor: '#e0e0e0',
        selectionBackground: '#3a3a3a',
        black: '#000000',
        red: '#ff5555',
        green: '#50fa7b',
        yellow: '#f1fa8c',
        blue: '#bd93f9',
        magenta: '#ff79c6',
        cyan: '#8be9fd',
        white: '#bfbfbf',
      },
    })

    const fitAddon = new FitAddon()
    term.loadAddon(fitAddon)

    term.open(terminalRef.current)
    fitAddon.fit()

    // Handle keyboard shortcuts
    term.attachCustomKeyEventHandler((event) => {
      // Ctrl+K for AI mode
      if (event.ctrlKey && event.key === 'k') {
        onOpenAI()
        return false
      }
      return true
    })

    xtermRef.current = term
    fitAddonRef.current = fitAddon

    // Initial message
    term.writeln('\x1b[1;34mAI Terminal Pro\x1b[0m - Press Ctrl+K for AI mode')
    term.writeln('')

    // Handle resize
    const handleResize = () => {
      fitAddon.fit()
    }

    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      term.dispose()
    }
  }, [onOpenAI])

  return (
    <div 
      ref={terminalRef} 
      className="h-full w-full p-2"
      style={{ backgroundColor: '#1a1a1a' }}
    />
  )
}
