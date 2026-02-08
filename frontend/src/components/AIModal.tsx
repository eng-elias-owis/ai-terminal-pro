import React, { useState, useCallback } from 'react'
import { X, Play, AlertTriangle, CheckCircle, AlertCircle } from 'lucide-react'

interface AIModalProps {
  onClose: () => void
  onExecute: (command: string) => void
}

export const AIModal: React.FC<AIModalProps> = ({ onClose, onExecute }) => {
  const [prompt, setPrompt] = useState('')
  const [generatedCommand, setGeneratedCommand] = useState('')
  const [riskLevel, setRiskLevel] = useState<'none' | 'low' | 'medium' | 'high' | 'critical'>('none')
  const [isGenerating, setIsGenerating] = useState(false)
  const [isEditing, setIsEditing] = useState(false)

  const handleGenerate = useCallback(async () => {
    if (!prompt.trim()) return
    
    setIsGenerating(true)
    
    // Simulate AI generation (will be replaced with actual API call)
    setTimeout(() => {
      setGeneratedCommand(`# AI would generate command for: ${prompt}`)
      setRiskLevel('low')
      setIsGenerating(false)
    }, 1000)
  }, [prompt])

  const handleExecute = useCallback(() => {
    if (generatedCommand) {
      onExecute(generatedCommand)
    }
  }, [generatedCommand, onExecute])

  const getRiskColor = () => {
    switch (riskLevel) {
      case 'none': return 'text-green-400'
      case 'low': return 'text-yellow-400'
      case 'medium': return 'text-orange-400'
      case 'high': return 'text-red-400'
      case 'critical': return 'text-red-600'
      default: return 'text-gray-400'
    }
  }

  const getRiskIcon = () => {
    switch (riskLevel) {
      case 'none': return <CheckCircle className="w-5 h-5 text-green-400" />
      case 'low':
      case 'medium': return <AlertCircle className="w-5 h-5 text-yellow-400" />
      case 'high':
      case 'critical': return <AlertTriangle className="w-5 h-5 text-red-400" />
      default: return null
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg w-full max-w-2xl m-4 p-6 shadow-2xl border border-gray-700">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-white">AI Command Generator</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        {/* Input */}
        <div className="mb-4">
          <label className="block text-sm text-gray-400 mb-2">
            What do you want to do? (Press Enter to generate)
          </label>
          <textarea
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="e.g., find all PDF files modified in the last 7 days"
            className="w-full h-20 bg-gray-900 border border-gray-700 rounded-lg p-3 text-white placeholder-gray-500 resize-none focus:outline-none focus:border-blue-500"
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault()
                handleGenerate()
              }
            }}
          />
        </div>

        {/* Generate Button */}
        {!generatedCommand && (
          <button
            onClick={handleGenerate}
            disabled={!prompt.trim() || isGenerating}
            className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 text-white rounded-lg font-medium transition-colors"
          >
            {isGenerating ? 'Generating...' : 'Generate Command'}
          </button>
        )}

        {/* Generated Command */}
        {generatedCommand && (
          <div className="space-y-4">
            {/* Risk Indicator */}
            <div className={`flex items-center gap-2 p-3 rounded-lg bg-gray-900 ${getRiskColor()}`}>
              {getRiskIcon()}
              <span className="text-sm font-medium capitalize">
                Risk Level: {riskLevel}
              </span>
            </div>

            {/* Command Display */}
            <div className="relative">
              {isEditing ? (
                <textarea
                  value={generatedCommand}
                  onChange={(e) => setGeneratedCommand(e.target.value)}
                  className="w-full h-20 bg-gray-900 border border-gray-700 rounded-lg p-3 text-white font-mono text-sm resize-none focus:outline-none focus:border-blue-500"
                />
              ) : (
                <pre className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white font-mono text-sm overflow-x-auto">
                  {generatedCommand}
                </pre>
              )}
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3">
              <button
                onClick={() => setIsEditing(!isEditing)}
                className="flex-1 py-2 px-4 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium transition-colors"
              >
                {isEditing ? 'Done Editing' : 'Edit'}
              </button>
              <button
                onClick={handleExecute}
                disabled={riskLevel === 'critical'}
                className="flex-1 py-2 px-4 bg-green-600 hover:bg-green-700 disabled:bg-gray-700 text-white rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
              >
                <Play className="w-4 h-4" />
                Execute
              </button>
            </div>
          </div>
        )}

        {/* Footer */}
        <div className="mt-4 pt-4 border-t border-gray-700 text-xs text-gray-500">
          Press Ctrl+Enter to execute, Esc to cancel
        </div>
      </div>
    </div>
  )
}
