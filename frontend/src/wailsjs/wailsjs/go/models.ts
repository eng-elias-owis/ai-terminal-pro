export namespace config {
	
	export class Settings {
	    litellm_endpoint: string;
	    model: string;
	    theme: string;
	    font_size: number;
	    font_family: string;
	    cursor_style: string;
	    ai_shortcut: string;
	    safety_mode: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.litellm_endpoint = source["litellm_endpoint"];
	        this.model = source["model"];
	        this.theme = source["theme"];
	        this.font_size = source["font_size"];
	        this.font_family = source["font_family"];
	        this.cursor_style = source["cursor_style"];
	        this.ai_shortcut = source["ai_shortcut"];
	        this.safety_mode = source["safety_mode"];
	    }
	}

}

