export namespace srvc {
	
	export class Developer {
	    id: number;
	    name: string;
	    email: string;
	    organization: string;
	
	    static createFrom(source: any = {}) {
	        return new Developer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.organization = source["organization"];
	    }
	}
	export class Environment {
	    id: number;
	    name: string;
	    group_name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.group_name = source["group_name"];
	        this.description = source["description"];
	    }
	}
	export class Packages {
	    id: number;
	    name: string;
	    title: string;
	    desc: string;
	    version: string;
	    repo: string;
	    work_dir: string;
	    developers: Developer;
	    supported_os: string[];
	    supported_arch: string[];
	
	    static createFrom(source: any = {}) {
	        return new Packages(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.desc = source["desc"];
	        this.version = source["version"];
	        this.repo = source["repo"];
	        this.work_dir = source["work_dir"];
	        this.developers = this.convertValues(source["developers"], Developer);
	        this.supported_os = source["supported_os"];
	        this.supported_arch = source["supported_arch"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Status {
	    active: boolean;
	    running: boolean;
	    startup: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active = source["active"];
	        this.running = source["running"];
	        this.startup = source["startup"];
	    }
	}
	export class Services {
	    id: number;
	    name: string;
	    title: string;
	    desc: string;
	    tags: string[];
	    port: number;
	    packages: Packages;
	    environments: Environment;
	    runtime_type: string;
	    runtime_config: Record<string, any>;
	    status: Status;
	
	    static createFrom(source: any = {}) {
	        return new Services(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.desc = source["desc"];
	        this.tags = source["tags"];
	        this.port = source["port"];
	        this.packages = this.convertValues(source["packages"], Packages);
	        this.environments = this.convertValues(source["environments"], Environment);
	        this.runtime_type = source["runtime_type"];
	        this.runtime_config = source["runtime_config"];
	        this.status = this.convertValues(source["status"], Status);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class VersionCheckResult {
	    Installed: boolean;
	    LocalVersion: string;
	    RemoteVersion: string;
	    UpdateAvailable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VersionCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Installed = source["Installed"];
	        this.LocalVersion = source["LocalVersion"];
	        this.RemoteVersion = source["RemoteVersion"];
	        this.UpdateAvailable = source["UpdateAvailable"];
	    }
	}

}

