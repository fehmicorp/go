import { useEffect, useState } from 'react';
import logo from './assets/images/logo.png';

interface Status {
    active: boolean;
    running: boolean;
    startup: boolean;
}

interface Package {
    id: number;
    name: string;
    title: string;
    desc: string;
    version: string;
    repo: string;
    workDir: string;
}

interface Environment {
    id: number;
    name: string;
    group_name: string;
    description: string;
}

interface Service {
    id: number;
    name: string;
    title: string;
    desc: string;
    tags: string[];
    port: number;
    package: Package;
    environment: Environment;
    runtime_type: string;
    status?: Status;
}

interface VersionCheckResult {
    Installed: boolean;
    LocalVersion: string;
    RemoteVersion: string;
    UpdateAvailable: boolean;
}

declare global {
    interface Window {
        go?: {
            main?: {
                App?: {
                    GetServiceList(): Promise<Service[]>;
                    ServiceAction(name: string, action: string): Promise<void>;
                    CheckServiceVersion(name: string): Promise<VersionCheckResult>;
                };
            };
        };
    }
}

function App() {
    const [services, setServices] = useState<Service[]>([]);
    const [serviceVersions, setServiceVersions] = useState<Record<string, VersionCheckResult>>({});
    const [loading, setLoading] = useState<boolean>(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    const fetchVersions = async (serviceList: Service[]) => {
        if (!window.go?.main?.App?.CheckServiceVersion) return;
        
        const versionMap: Record<string, VersionCheckResult> = {};
        for (const srv of serviceList) {
            try {
                const vRes = await window.go.main.App.CheckServiceVersion(srv.name);
                versionMap[srv.name] = vRes;
            } catch (err) {
                console.error(`Failed to fetch version for ${srv.name}:`, err);
            }
        }
        setServiceVersions(versionMap);
    };

    const fetchServices = async () => {
        try {
            if (window.go?.main?.App?.GetServiceList) {
                const list = await window.go.main.App.GetServiceList();
                const validList = list || [];
                setServices(validList);
                await fetchVersions(validList);
            }
        } catch (err) {
            console.error("Failed to fetch service list:", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchServices();
        const interval = setInterval(fetchServices, 5000);
        return () => clearInterval(interval);
    }, []);

    const handleAction = async (serviceName: string, actionType: string) => {
        try {
            setActionLoading(`${serviceName}-${actionType}`);
            if (window.go?.main?.App?.ServiceAction) {
                await window.go.main.App.ServiceAction(serviceName, actionType);
                await fetchServices();
            }
        } catch (err) {
            console.error(`Failed to execute ${actionType} for ${serviceName}:`, err);
        } finally {
            setActionLoading(null);
        }
    };

    return (
        <div id="App" className="flex flex-col items-center justify-start min-h-screen p-6" style={{ backgroundColor: 'var(--bg-base)', color: 'var(--text-main)' }}>
            <div className="p-6 rounded-2xl shadow-2xl flex flex-col w-full max-w-xl backdrop-blur-md" style={{ backgroundColor: 'var(--bg-card)', border: '1px solid var(--border-card)' }}>
                
                {/* Header Section */}
                <div className='flex flex-row items-center justify-between pb-4 mb-4' style={{ borderBottom: '1px solid var(--border-color)' }}>
                    <div className='flex flex-row items-center gap-3'>
                        <img src={logo} id="logo" alt="logo" className="w-9 h-9 object-contain drop-shadow-md" />
                        <div>
                            <h1 className='text-lg font-bold tracking-wide'>Fehmi Cloud Connector</h1>
                            <p className='text-xs font-light' style={{ color: 'var(--text-muted)' }}>Service Manager & Monitoring</p>
                        </div>
                    </div>
                    <button 
                        onClick={fetchServices}
                        className="px-3 py-1.5 text-xs transition rounded-lg font-medium"
                        style={{ backgroundColor: 'var(--bg-elevated)', border: '1px solid var(--border-color)', color: 'var(--text-main)' }}>
                        Refresh
                    </button>
                </div>

                {/* Services List Container */}
                <div className="flex flex-col w-full gap-3">
                    {loading ? (
                        <div className="text-center py-8 text-sm" style={{ color: 'var(--text-muted)' }}>Loading services...</div>
                    ) : services.length === 0 ? (
                        <div className="text-center py-8 text-sm" style={{ color: 'var(--text-muted)' }}>No services registered.</div>
                    ) : (
                        services.map((service) => {
                            const isRunning = service.status?.running ?? false;
                            const isActive = service.status?.active ?? false;
                            const isStartupAuto = service.status?.startup ?? false;
                            const verInfo = serviceVersions[service.name];

                            return (
                                <div 
                                    key={service.id || service.name} 
                                    className="flex flex-col rounded-xl p-4 transition gap-3"
                                    style={{ backgroundColor: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                                    
                                    <div className="flex flex-row items-center justify-between">
                                        <div className="flex flex-col">
                                            <span className="text-sm font-semibold" style={{ color: 'var(--text-main)' }}>{service.title || service.name}</span>
                                            <div className="flex flex-wrap items-center gap-2 mt-0.5">
                                                <span className="text-[11px] px-2 py-0.5 rounded font-mono" style={{ backgroundColor: 'var(--bg-surface)', color: 'var(--text-muted)' }}>
                                                    Port: {service.port}
                                                </span>
                                                <span className="text-[11px] px-2 py-0.5 rounded font-mono" style={{ backgroundColor: 'var(--bg-surface)', color: 'var(--text-muted)' }}>
                                                    Runtime: {service.runtime_type}
                                                </span>
                                                <span className="text-[11px] px-2 py-0.5 rounded font-mono" style={{ backgroundColor: 'var(--bg-surface)', color: 'var(--text-muted)' }}>
                                                    Installed: {verInfo?.Installed ? (verInfo.LocalVersion || 'Yes') : 'Not Installed'}
                                                </span>
                                                <span className="text-[11px] px-2 py-0.5 rounded font-mono" style={{ backgroundColor: 'var(--bg-surface)', color: 'var(--text-muted)' }}>
                                                    Listed: {service.package?.version || 'N/A'}
                                                </span>
                                            </div>
                                        </div>

                                        {/* Status Indicator Badge */}
                                        <div className="flex items-center gap-2">
                                            <span className={`h-2.5 w-2.5 rounded-full ${isRunning ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.6)]' : 'bg-rose-400'}`}></span>
                                            <span className="text-xs font-medium" style={{ color: 'var(--text-main)' }}>
                                                {isRunning ? 'Running' : 'Stopped'}
                                            </span>
                                        </div>
                                    </div>

                                    {/* Action Controls */}
                                    <div className="flex flex-row items-center justify-between pt-2 text-xs" style={{ borderTop: '1px solid var(--border-subtle)' }}>
                                        <div className="flex items-center gap-3 text-[11px]" style={{ color: 'var(--text-muted)' }}>
                                            <span>Active: {isActive ? 'Yes' : 'No'}</span>
                                            <span>•</span>
                                            <span>Startup: {isStartupAuto ? 'Auto' : 'Manual'}</span>
                                        </div>

                                        <div className="flex items-center gap-2">
                                            {isRunning ? (
                                                <>
                                                    <button 
                                                        disabled={actionLoading === `${service.name}-restart`}
                                                        onClick={() => handleAction(service.name, 'restart')}
                                                        className="px-3 py-1 bg-amber-500/25 hover:bg-amber-500/35 text-amber-300 border border-amber-500/40 rounded-lg transition font-medium disabled:opacity-50">
                                                        {actionLoading === `${service.name}-restart` ? 'Processing...' : 'Restart'}
                                                    </button>
                                                    <button 
                                                        disabled={actionLoading === `${service.name}-stop`}
                                                        onClick={() => handleAction(service.name, 'stop')}
                                                        className="px-3 py-1 bg-rose-500/25 hover:bg-rose-500/35 text-rose-300 border border-rose-500/40 rounded-lg transition font-medium disabled:opacity-50">
                                                        {actionLoading === `${service.name}-stop` ? 'Processing...' : 'Stop'}
                                                    </button>
                                                </>
                                            ) : (
                                                <button 
                                                    disabled={actionLoading === `${service.name}-start`}
                                                    onClick={() => handleAction(service.name, 'start')}
                                                    className="px-3 py-1 bg-emerald-500/25 hover:bg-emerald-500/35 text-emerald-300 border border-emerald-500/40 rounded-lg transition font-medium disabled:opacity-50">
                                                    {actionLoading === `${service.name}-start` ? 'Processing...' : 'Start'}
                                                </button>
                                            )}
                                        </div>
                                    </div>

                                </div>
                            );
                        })
                    )}
                </div>

            </div>
        </div>
    );
}

export default App;