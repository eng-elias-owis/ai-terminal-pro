export {};

declare global {
  interface Window {
    runtime: {
      EventsOn: (eventName: string, callback: (...data: any) => void) => () => void;
      EventsOff: (eventName: string) => void;
      EventsEmit: (eventName: string, ...data: any) => void;
    };
    go: {
      main: {
        App: {
          WriteToTerminal: (data: string) => Promise<void>;
          ResizeTerminal: (rows: number, cols: number) => Promise<void>;
          GenerateCommand: (description: string) => Promise<Record<string, any>>;
          ValidateCommand: (command: string) => Promise<Record<string, any>>;
          GetSettings: () => Promise<any>;
          SaveSettings: (settings: any) => Promise<void>;
          Greet: (name: string) => Promise<string>;
        };
      };
    };
  }
}
