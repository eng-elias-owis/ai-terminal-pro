import React, { useState, useCallback } from 'react'
import { Terminal } from './components/Terminal'
import { AIModal } from './components/AIModal'
import { Settings } from './components/Settings'
import { Settings as SettingsIcon, Terminal as TerminalIcon } from 'lucide-react'

function App() {
  const [showSettings, setShowSettings] = useState(false)
  const [showAIModal, setShowAIModal] = useState(false)

  const handleOpenAI = useCallback(() => {
    setShowAIModal(true)
  }, [])

  const handleCloseAI = useCallback(() => {
    setShowAIModal(false)
  }, [])

  const handleExecuteCommand = useCallback((command: string) => {
    // Send command to terminal
    setShowAIModal(false)
  }, [])

  return (
    <div className="h-screen w-screen bg-gray-900 text-white flex flex-col">
      {/* Header */}
      <div className="h-12 bg-gray-800 border-b border-gray-700 flex items-center px-4 justify-between">
        <div className="flex items-center gap-2">
          <TerminalIcon className="w-5 h-5 text-blue-400" />
          <span className="font-semibold">AI Terminal Pro</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowSettings(true)}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
            title="Settings"
          >
            <SettingsIcon className="w-5 h-5" />
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 relative">
        <Terminal onOpenAI={handleOpenAI} />
        
        {showAIModal && (
          <AIModal
            onClose={handleCloseAI}
            onExecute={handleExecuteCommand}
          />
        )}

        {showSettings && (
          <Settings onClose={() => setShowSettings(false)} />
        )}
      </div>

      {/* Status Bar */}
      <div className="h-6 bg-gray-800 border-t border-gray-700 flex items-center px-4 text-xs text-gray-400">
        <span>Ready</span>
        <span className="ml-auto">Press Ctrl+K for AI mode</span>
      </div>
    </div>
  )
}

export default App
