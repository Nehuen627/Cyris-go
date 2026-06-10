export namespace structs {
	
	export class CPUResult {
	    Meets: boolean;
	    MeetsRecommended: boolean;
	    Found: boolean;
	    UserCPU: string;
	    RequiredCPU: string;
	
	    static createFrom(source: any = {}) {
	        return new CPUResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Meets = source["Meets"];
	        this.MeetsRecommended = source["MeetsRecommended"];
	        this.Found = source["Found"];
	        this.UserCPU = source["UserCPU"];
	        this.RequiredCPU = source["RequiredCPU"];
	    }
	}
	export class GPUResult {
	    Meets: boolean;
	    MeetsRecommended: boolean;
	    Found: boolean;
	    UserGPU: string;
	    RequiredGPU: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Meets = source["Meets"];
	        this.MeetsRecommended = source["MeetsRecommended"];
	        this.Found = source["Found"];
	        this.UserGPU = source["UserGPU"];
	        this.RequiredGPU = source["RequiredGPU"];
	    }
	}
	export class Requirements {
	    minimum: string;
	    recommended: string;
	
	    static createFrom(source: any = {}) {
	        return new Requirements(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minimum = source["minimum"];
	        this.recommended = source["recommended"];
	    }
	}
	export class GameData {
	    pc_requirements: Requirements;
	    mac_requirements: Requirements;
	    linux_requirements: Requirements;
	
	    static createFrom(source: any = {}) {
	        return new GameData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pc_requirements = this.convertValues(source["pc_requirements"], Requirements);
	        this.mac_requirements = this.convertValues(source["mac_requirements"], Requirements);
	        this.linux_requirements = this.convertValues(source["linux_requirements"], Requirements);
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
	
	export class RequirementsResult {
	    CPUCores: boolean;
	    RAMTotal: boolean;
	    DiskFree: boolean;
	    GPU: GPUResult;
	    CPU: CPUResult;
	    MeetsMinimum: boolean;
	    MeetsRecommended: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RequirementsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CPUCores = source["CPUCores"];
	        this.RAMTotal = source["RAMTotal"];
	        this.DiskFree = source["DiskFree"];
	        this.GPU = this.convertValues(source["GPU"], GPUResult);
	        this.CPU = this.convertValues(source["CPU"], CPUResult);
	        this.MeetsMinimum = source["MeetsMinimum"];
	        this.MeetsRecommended = source["MeetsRecommended"];
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
	export class SteamApp {
	    appid: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SteamApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appid = source["appid"];
	        this.name = source["name"];
	    }
	}
	export class SystemSpecs {
	    cpu_name: string;
	    cpu_cores: number;
	    ram_total_mb: number;
	    disk_free_gb: number;
	    gpu_name: string;
	    gpu_vendor: string;
	    os: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemSpecs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu_name = source["cpu_name"];
	        this.cpu_cores = source["cpu_cores"];
	        this.ram_total_mb = source["ram_total_mb"];
	        this.disk_free_gb = source["disk_free_gb"];
	        this.gpu_name = source["gpu_name"];
	        this.gpu_vendor = source["gpu_vendor"];
	        this.os = source["os"];
	    }
	}

}

