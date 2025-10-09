# Known Fixes and Patches

## Stack Overflow in RequestCtxSpanExtractor (2025-10-09)

### Problem
Infinite recursion causing stack overflow when `trace.SpanFromContext` is called in applications using `host-fasthttp` with tracing enabled.

### Root Cause
1. `FasthttpHost` initialization registers `RequestCtxSpanExtractor` as global span extractor
2. When `trace.SpanFromContext(ctx)` is called, it triggers `globalSpanExtractor.Extract(ctx)`
3. `RequestCtxSpanExtractor.Extract()` calls `trace.SpanFromContext(ctx)` as fallback
4. This creates infinite recursion: `Extract → SpanFromContext → Extract → ...`

### Recursion Flow
```textplain
1. FasthttpHost.init() sets globalSpanExtractor = RequestCtxSpanExtractor(0)
2. Application calls trace.SpanFromContext(ctx)
3. trace.SpanFromContext calls globalSpanExtractor.Extract(ctx)
4. RequestCtxSpanExtractor.Extract checks requestutil.ExtractSpan(rx)
5. If span not found, calls trace.SpanFromContext(ctx) again
6. Loop back to step 3 → Stack Overflow
```

### Solution
Modified [internal/tracingutil/requestCtxSpanExtractor.go](internal/tracingutil/requestCtxSpanExtractor.go#L29):

```go
// Before (causes infinite recursion)
func (RequestCtxSpanExtractor) Extract(ctx context.Context) *trace.SeveritySpan {
    if rx, ok := ctx.(*http.RequestCtx); ok {
        reply := requestutil.ExtractSpan(rx)
        span, ok := reply.(*trace.SeveritySpan)
        if ok {
            return span
        }
    }
    return trace.SpanFromContext(ctx)  // ← Infinite recursion
}

// After (fixed)
func (RequestCtxSpanExtractor) Extract(ctx context.Context) *trace.SeveritySpan {
    if rx, ok := ctx.(*http.RequestCtx); ok {
        reply := requestutil.ExtractSpan(rx)
        span, ok := reply.(*trace.SeveritySpan)
        if ok {
            return span
        }
    }
    // Return nil instead of calling trace.SpanFromContext to avoid infinite recursion
    return nil
}
```

### Impact
- ✅ Fixes stack overflow in production applications
- ✅ Maintains tracing functionality
- ✅ Follows Go's zero-value convention
- ✅ Minimal change (1 line)
- ✅ No breaking changes to public API

### Branch
`fix/stack-overflow-span-extractor`

### Commit
`eb624f9` - fix(trace): prevent infinite recursion in RequestCtxSpanExtractor

### Related Files
- [internal/tracingutil/requestCtxSpanExtractor.go:29](internal/tracingutil/requestCtxSpanExtractor.go#L29) (fix applied)
- [internal/requestTracerService.go:51](internal/requestTracerService.go#L51) (sets global extractor)
- [internal/def.go:29](internal/def.go#L29) (defines defaultSpanExtractor)

### Testing
- ✅ All existing tests pass (except `TestStartup_UseTracing` which requires external Jaeger service)
- ✅ Code builds successfully
- ✅ No new compilation errors
- ✅ Tracingutil package tests pass

---
**Fixed**: 2025-10-09
**Severity**: Critical (Stack Overflow)
**Affected Versions**: All versions with tracing enabled
