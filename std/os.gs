export {
    chdir: func(path) {
        __chdir(path);
    },
    exit: func(code) {
        __exit(code);
    },
    getEnv: func(key) {
        return __getenv(key);
    },
    setEnv: func(key, val) {
        return __setenv(key, val);
    },
    args: func() {
        return __args();
    },
	getegid: func() {
        return __getegid();
    },
	geteuid: func() {
        return __geteuid();
    },
	getgid: func() {
        return __getgid();
    },
	getpid: func() {
        return __getpid();
    },
    getppid: func() {
        return __getppid();
    },
	getuid: func() {
        return __getuid();
    },
    exec: func(cmd,...args) {
        return __exec(cmd, args);
    },
}