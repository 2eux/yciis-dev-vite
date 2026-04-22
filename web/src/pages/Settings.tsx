import { useState } from 'react'
import { Building, Users, Bell, Lock, Palette, Globe, Database, Save } from 'lucide-react'

const tabs = [
  { id: 'general', label: 'General', icon: Building },
  { id: 'users', label: 'Users & Roles', icon: Users },
  { id: 'notifications', label: 'Notifications', icon: Bell },
  { id: 'security', label: 'Security', icon: Lock },
  { id: 'appearance', label: 'Appearance', icon: Palette },
  { id: 'integration', label: 'Integrations', icon: Globe },
  { id: 'backup', label: 'Backup', icon: Database },
]

export default function Settings() {
  const [activeTab, setActiveTab] = useState('general')

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Settings</h1>
        <p className="text-slate-500">Configure system settings and preferences</p>
      </div>

      <div className="flex gap-6">
        <div className="w-64 shrink-0">
          <nav className="space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`w-full flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-colors ${
                  activeTab === tab.id
                    ? 'bg-blue-50 text-blue-600'
                    : 'text-slate-600 hover:bg-slate-50'
                }`}
              >
                <tab.icon className="w-5 h-5" />
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        <div className="flex-1">
          <div className="card">
            {activeTab === 'general' && (
              <div className="space-y-6">
                <h2 className="text-lg font-semibold text-slate-900">School Information</h2>
                
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      School Name
                    </label>
                    <input
                      type="text"
                      defaultValue="Edusys Pro Academy"
                      className="input"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      School Code
                    </label>
                    <input
                      type="text"
                      defaultValue="EDU001"
                      className="input"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      Email
                    </label>
                    <input
                      type="email"
                      defaultValue="info@edusys.edu"
                      className="input"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      Phone
                    </label>
                    <input
                      type="text"
                      defaultValue="+62 21 1234 5678"
                      className="input"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      Timezone
                    </label>
                    <select className="input">
                      <option>Asia/Jakarta (WIB)</option>
                      <option>Asia/Makassar (WITA)</option>
                      <option>Asia/Jayapura (WIT)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">
                      Currency
                    </label>
                    <select className="input">
                      <option>IDR (Indonesian Rupiah)</option>
                      <option>USD (US Dollar)</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">
                    Address
                  </label>
                  <textarea
                    rows={3}
                    defaultValue="Jl. Education No. 123, Jakarta Selatan"
                    className="input"
                  />
                </div>

                <div className="flex justify-end">
                  <button className="btn btn-primary">
                    <Save className="w-4 h-4" />
                    Save Changes
                  </button>
                </div>
              </div>
            )}

            {activeTab === 'security' && (
              <div className="space-y-6">
                <h2 className="text-lg font-semibold text-slate-900">Security Settings</h2>
                
                <div className="space-y-4">
                  <label className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
                    <div>
                      <p className="font-medium text-slate-900">Two-Factor Authentication</p>
                      <p className="text-sm text-slate-500">Require 2FA for all admin accounts</p>
                    </div>
                    <input type="checkbox" className="toggle" />
                  </label>
                  
                  <label className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
                    <div>
                      <p className="font-medium text-slate-900">Password Expiry</p>
                      <p className="text-sm text-slate-500">Force password change every 90 days</p>
                    </div>
                    <input type="checkbox" className="toggle" />
                  </label>
                  
                  <label className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
                    <div>
                      <p className="font-medium text-slate-900">Session Timeout</p>
                      <p className="text-sm text-slate-500">Auto logout after 30 minutes of inactivity</p>
                    </div>
                    <input type="checkbox" defaultChecked className="toggle" />
                  </label>
                </div>
              </div>
            )}

            {activeTab !== 'general' && activeTab !== 'security' && (
              <div className="text-center py-12">
                <p className="text-slate-500">{tabs.find(t => t.id === activeTab)?.label} settings coming soon...</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}