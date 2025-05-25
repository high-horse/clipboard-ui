export namespace clip {
	
	export class CopiedContent {
	    key: number;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new CopiedContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}

}

