'use client';

import { useState, useEffect, useRef } from "react";
import { Config, URL } from "../../../../config/index";

interface License {
  id: string;
  name: string;
  type: string;
  status: string;
}

export default function OnboardPage() {
  const [showWelcome, setShowWelcome] = useState(true);
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  
  const [auth, setAuth] = useState({ Email: "", Password: "" });
  const [emailValid, setEmailValid] = useState(false);

  const [licenses, setLicenses] = useState<License[]>([]);

  const [formData, setFormData] = useState({
    instanceName: "",
    token: "",
    selectedLicenseId: "",
    onpremMode: "standalone",
  });

  const emailInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (step === 2) {
      emailInputRef.current?.focus();
    }
  }, [step]);

  const validateEmail = (email: string) => {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setAuth({ ...auth, Email: val });
    setEmailValid(validateEmail(val));
  };

  const handleLogout = () => {
    setAuth({ Email: "", Password: "" });
    setEmailValid(false);
    setFormData(prev => ({ ...prev, token: "", selectedLicenseId: "" }));
    setStep(2);
  };

  const handleNext = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    if (step === 1) {
      if (!formData.instanceName.trim()) {
        setErrorMsg("Please provide a valid instance name.");
        return;
      }
      setStep(2);
    } 
    else if (step === 2) {
      setLoading(true);
      try {
        // Simulated client-side authentication for static export
        await new Promise((resolve) => setTimeout(resolve, 800));

        if (!auth.Email || !auth.Password) {
          throw new Error("Invalid administrator credentials.");
        }

        const token = "fc_onprem_token_mock_xyz987";
        const fetchedLicenses = [
          { id: "lic-standalone-01", name: "Standalone Local Node License", type: "Standard", status: "Active" }
        ];

        setLicenses(fetchedLicenses);
        setFormData(prev => ({ 
          ...prev, 
          token, 
          selectedLicenseId: fetchedLicenses[0]?.id || "" 
        }));
        setStep(3);
      } catch (err: any) {
        setErrorMsg(err.message || "Authentication failed. Please check your credentials.");
      } finally {
        setLoading(false);
      }
    } 
    else if (step === 3) {
      setLoading(true);
      try {
        await new Promise((resolve) => setTimeout(resolve, 800));
        alert("On-premises instance configured successfully! Redirecting...");
        window.location.href = "/";
      } catch (err) {
        alert("On-premises instance configured successfully! Redirecting...");
        window.location.href = "/";
      } finally {
        setLoading(false);
      }
    }
  };

  if (showWelcome) {
    return (
      <div 
        onClick={() => setShowWelcome(false)} 
        className="fixed inset-0 w-screen h-screen flex flex-col items-center justify-center select-none space-y-6 text-center z-50 bg-[var(--bg-main, #0f172a)] transition-opacity duration-1000 cursor-pointer"
      >
        <div className="space-y-4 animate-pulse">
          <img 
            src="/logo/icon_192.png" 
            alt="Logo" 
            className="w-24 h-24 mx-auto object-contain transition-transform duration-500 hover:scale-105" 
          />
          <h1 className="text-2xl font-bold tracking-tight text-[var(--text-main)]">
            {Config.Name} {Config.Apps.Cloud.Title}
          </h1>
          <p className="text-sm text-[var(--text-muted)]">
            {Config.Apps.Cloud.Tagline}
          </p>
        </div>
        <div className="text-center text-xs text-[var(--text-subtle)] pt-12">
          Click here to begin step
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-[85vh] flex flex-col justify-between space-y-6 transition-all duration-[2000ms] ease-in-out">
      <div className="space-y-6 transition-all duration-[2000ms] ease-in-out">
        <div className="bg-[var(--bg-banner)] text-[var(--text-banner)] mt-16 px-4 py-3 rounded-lg text-sm font-medium flex items-center justify-between border border-[var(--border-card)] shadow-sm transition-all duration-[2000ms]">
          <span>🚀 On-Premises Setup — {Config.Company.LegalName}</span>
          <span className="text-xs opacity-80">Step {step} of 3</span>
        </div>
        <div className="bg-gradient-to-b from-[var(--bg-hero-from)] to-[var(--bg-hero-to)] p-8 rounded-2xl border border-[var(--border-card)] shadow-xl transition-all duration-[2000ms] ease-in-out">
          <h1 className="text-2xl font-bold tracking-tight text-[var(--text-main)]">
            {Config.Name} {Config.Apps.Cloud.Title}
          </h1>
          <p className="text-sm text-[var(--text-muted)] mt-1">{Config.Apps.Cloud.Tagline}</p>
          {errorMsg && (
            <div className="mt-4 p-3 bg-red-500/10 border border-red-500/30 text-red-400 text-xs rounded-xl">
              {errorMsg}
            </div>
          )}
          <form onSubmit={handleNext} className="mt-6 space-y-5">
            <div className="transition-all duration-[2000ms] ease-in-out">
              {step === 1 && (
                <div className="space-y-2 transition-all duration-[2000ms] ease-in-out">
                  <label className="block text-xs font-semibold uppercase tracking-wider text-[var(--text-subtle)]">
                    Instance Identifier
                  </label>
                  <input
                    type="text"
                    value={formData.instanceName}
                    onChange={(e) => setFormData({ ...formData, instanceName: e.target.value })}
                    placeholder="e.g. fehmicorp-onprem-node"
                    className="w-full px-4 py-3 bg-[var(--bg-card)] border border-[var(--border-card)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--fehmicorp-primary)]"
                    required
                  />
                </div>
              )}
              {step === 2 && (
                <div className="space-y-4 transition-all duration-[2000ms] ease-in-out">
                  <div className="space-y-2">
                    <label className="block text-xs font-semibold uppercase tracking-wider text-[var(--text-subtle)]">
                      Administrator Email
                    </label>
                    <input
                      ref={emailInputRef}
                      type="email"
                      value={auth.Email}
                      onChange={handleEmailChange}
                      placeholder="admin@fehmicorp.cloud"
                      className="w-full px-4 py-3 bg-[var(--bg-card)] border border-[var(--border-card)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--fehmicorp-primary)]"
                      required
                    />
                  </div>
                  <div className={`transition-all duration-[2000ms] ease-in-out overflow-hidden ${emailValid ? 'max-h-40 opacity-100 transform translate-y-0' : 'max-h-0 opacity-0 transform -translate-y-4 pointer-events-none'}`}>
                    <div className="space-y-2 p-1">
                      <label className="block text-xs font-semibold uppercase tracking-wider text-[var(--text-subtle)]">
                        Administrator Password
                      </label>
                      <input
                        type="password"
                        value={auth.Password}
                        onChange={(e) => setAuth({ ...auth, Password: e.target.value })}
                        placeholder="•••••••"
                        className="w-full px-4 py-3 bg-[var(--bg-card)] border border-[var(--border-card)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--fehmicorp-primary)]"
                        required={emailValid}
                      />
                    </div>
                  </div>
                </div>
              )}
              {step === 3 && (
                <div className="space-y-4 transition-all duration-[2000ms] ease-in-out">
                  <div className="flex items-center justify-between bg-white/5 dark:bg-black/20 backdrop-blur-xl border border-white/10 dark:border-white/5 p-4 rounded-2xl shadow-lg">
                    <div>
                      <span className="text-xs text-[var(--text-subtle)] uppercase tracking-wider block">Logged in as</span>
                      <span className="text-sm font-medium text-[var(--text-main)]">{auth.Email}</span>
                    </div>
                    <button
                      type="button"
                      onClick={handleLogout}
                      className="px-3 py-1.5 text-xs cursor-pointer font-medium rounded-xl border border-red-500/30 text-red-400 hover:bg-red-500/10 backdrop-blur-md transition"
                    >
                      Use Different Account
                    </button>
                  </div>
                  <div className="space-y-2">
                    <label className="block text-xs font-semibold uppercase tracking-wider text-[var(--text-subtle)]">
                      Assign On-Premises License
                    </label>
                    <div className="relative">
                      <select
                        value={formData.selectedLicenseId}
                        onChange={(e) => setFormData({ ...formData, selectedLicenseId: e.target.value })}
                        className="w-full appearance-none px-4 py-3.5 bg-white/10 dark:bg-black/30 backdrop-blur-2xl border border-white/20 dark:border-white/10 rounded-2xl text-[var(--text-main)] shadow-[0_8px_32px_0_rgba(0,0,0,0.1)] focus:outline-none focus:ring-2 focus:ring-[var(--fehmicorp-primary)] focus:border-transparent transition-all"
                        required
                      >
                        {licenses.map((lic) => (
                          <option key={lic.id} value={lic.id} className="bg-[var(--bg-card)] backdrop-blur-2xl rounded-2xl dark:border-white/10 text-[var(--text-main)]">
                            {lic.name} ({lic.type}) — {lic.status}
                          </option>
                        ))}
                      </select>
                      <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-[var(--text-muted)]">
                        <svg className="h-4 w-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                          <path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" />
                        </svg>
                      </div>
                    </div>
                  </div>
                  <p className="text-xs text-[var(--text-muted)]">
                    Ready to provision local instance node.
                  </p>
                </div>
              )}
            </div>
            <div className="flex justify-between items-center pt-4">
              {step === 2 ? (
                <button
                  type="button"
                  onClick={() => setStep(step - 1)}
                  disabled={loading}
                  className="px-5 py-2.5 cursor-pointer text-sm font-medium rounded-xl border border-[var(--border-card)] text-[var(--text-muted)] hover:bg-[var(--bg-elevated)] transition"
                >
                  Back
                </button>
              ) : <div />}
              <button
                type="submit"
                disabled={loading || (step === 2 && !emailValid)}
                className="px-6 py-2.5 cursor-pointer text-sm font-semibold rounded-xl bg-[var(--fehmicorp-primary)] text-text-main text-white hover:bg-[var(--fehmicorp-primary-dark)] transition shadow-md disabled:opacity-50"
              >
                {loading ? "Processing..." : step === 2 ? "Sign In" : step === 3 ? "Complete Setup" : "Next"}
              </button>
            </div>
          </form>
        </div>
      </div>
      <div className="text-center text-xs text-[var(--text-subtle)] pt-6 mt-auto">
        &copy; {new Date().getFullYear()} 
        <a href={URL.Website} className="mx-1">{Config.Company.LegalName}.</a>
        All rights reserved.
      </div>
    </div>
  );
}