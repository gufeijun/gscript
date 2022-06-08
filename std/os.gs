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
}