import React, { useState, useEffect } from 'react'
import { X, Save, Eye, EyeOff } from 'lucide-react'

interface SettingsProps {
  onClose: () => void
}

export const Settings: React.FC<SettingsProps> = ({ onClose }) => {
  const [endpoint, setEndpoint] = useState('')
  const [virtualKey, setVirtualKey] = useState('')
  const [showKey, setShowKey] = useState(false)
  const [theme, setTheme] = useState('dark')
  const [fontSize, setFontSize] = useState(14)
  const [safetyMode, setSafetyMode] = useState('normal')

  useEffect(() => {
    // Load settings from backend (placeholder)
    // Will be implemented with actual backend integration
  }, [])

  const handleSave = () => {
    // Save settings to backend
    console.log('Saving settings:', { endpoint, virtualKey, theme, fontSize, safetyMode })
    onClose()
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg w-full max-w-lg m-4 p-6 shadow-2xl border border-gray-700">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-white">Settings</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        <div className="space-y-6">
          {/* LiteLLM Endpoint */}
          <div>
            <label className="block text-sm text-gray-400 mb-2">
              LiteLLM Endpoint URL
            </label>
            <input
              type="text"
              value={endpoint}
              onChange={(e) => setEndpoint(e.target.value)}
              placeholder="https://your-space.hf.space"
              className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
            />
          </div>

          {/* Virtual Key */}
          <div>
            <label className="block text-sm text-gray-400 mb-2">
              LiteLLM Virtual Key
            </label>
            <div className="relative">
              <input
                type={showKey ? 'text' : 'password'}
                value={virtualKey}
                onChange={(e) => setVirtualKey(e.target.value)}
                placeholder="sk-litellm-xxxxx"
                className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 pr-12 text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
              />
              <button
                onClick={() => setShowKey(!showKey)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white"
              >
                {showKey ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
            <p className="mt-1 text-xs text-gray-500">
              Your virtual key is stored securely in the OS keyring
            </p>
          </div>

          {/* Appearance */}
          <div className="border-t border-gray-700 pt-4">
            <h3 className="text-sm font-medium text-white mb-3">Appearance</h3>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm text-gray-400 mb-2">Theme</label>
                <select
                  value={theme}
                  onChange={(e) => setTheme(e.target.value)}
                  className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:border-blue-500"
                >
                  <option value="dark">Dark</option>
                  <option value="light">Light</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-gray-400 mb-2">Font Size</label>
                <input
                  type="number"
                  value={fontSize}
                  onChange={(e) => setFontSize(Number(e.target.value))}
                  min={10}
                  max={24}
                  className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:border-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Security */}
          <div className="border-t border-gray-700 pt-4">
            <h3 className="text-sm font-medium text-white mb-3">Security</h3>
            
            <div>
              <label className="block text-sm text-gray-400 mb-2">Safety Mode</label>
              <select
                value={safetyMode}
                onChange={(e) => setSafetyMode(e.target.value)}
                className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:border-blue-500"
              >
                <option value="strict">Strict - Block all risky commands</option>
                <option value="normal">Normal - Confirm risky commands</option>
                <option value="off">Off - No safety checks</option>
              </select>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-6 pt-4 border-t border-gray-700 flex justify-end">
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
          >
            <Save className="w-4 h-4" />
            Save Settings
          </button>
        </div>
      </div>
    </div>
  )
}
