package state

import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_absindex
// lua-5.3.4/src/lapi.c#lua_absindex()
func (self *luaState) AbsIndex(idx int) int {
	if idx > 0 || _isPseudo(idx) {
		return idx
	}
	return self.stack.absIndex(idx)
}

/* test for pseudo index */
func _isPseudo(i int) bool {
	return i <= LUA_REGISTRYINDEX
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_checkstack
// lua-5.3.4/src/lapi.c#lua_checkstack()
func (self *luaState) CheckStack(n int) bool {
	return self.stack.check(n)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gettop
// lua-5.3.4/src/lapi.c#lua_gettop()
func (self *luaState) GetTop() int {
	return self.stack.sp
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_settop
// lua-5.3.4/src/lapi.c#lua_settop()
func (self *luaState) SetTop(idx int) {
	if idx < 0 {
		idx = self.stack.absIndex(idx)
	}

	n := self.stack.sp - idx
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}

// [-n, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pop
// lua-5.3.4/src/lua.h#lua_pop()
func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvalue
// lua-5.3.4/src/lapi.c#lua_pushvalue()
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_copy
// lua-5.3.4/src/lapi.c#lua_copy()
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_insert
// lua-5.3.4/src/lua.h#lua_insert()
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_remove
// lua-5.3.4/src/lua.h#lua_remove()
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_replace
// lua-5.3.4/src/lua.h#lua_replace()
func (self *luaState) Replace(idx int) {
	self.Copy(-1, idx)
	self.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rotate
// lua-5.3.4/src/lapi.c#lua_rotate()
func (self *luaState) Rotate(idx, n int) {
	stack := self.stack
	t := stack.sp - 1        /* end of stack segment being rotated */
	p := stack.absIndex(idx) /* start of segment */
	p -= 1
	var m int /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	stack.reverse(m+1, t) /* reverse the suffix */
	stack.reverse(p, t)   /* reverse the entire segment */
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_xmove
// lua-5.3.4/src/lapi.c#lua_rotate()
func (self *luaState) XMove(to LuaState, n int) {
	lsFrom := self
	lsTo := to.(*luaState)

	elems := lsFrom.stack.popN(n)
	lsTo.stack.pushN(elems)
}
