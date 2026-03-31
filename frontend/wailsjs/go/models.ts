export namespace domain {
	
	export class Result___map_string_interface____ {
	    success: boolean;
	    data?: any[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result___map_string_interface____(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class Result_bool_ {
	    success: boolean;
	    data?: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result_bool_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class Result_map_string_interface____ {
	    success: boolean;
	    data?: Record<string, any>;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result_map_string_interface____(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class Result_string_ {
	    success: boolean;
	    data?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}

}

