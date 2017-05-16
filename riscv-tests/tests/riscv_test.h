#ifndef _ENV_PICORV32_TEST_H
#define _ENV_PICORV32_TEST_H

#ifndef TEST_FUNC_NAME
#  define TEST_FUNC_NAME mytest
#  define TEST_FUNC_TXT "mytest"
#  define TEST_FUNC_RET mytest_ret
#endif

#define RVTEST_RV32U
#define TESTNUM x28

#define RVTEST_CODE_BEGIN												\
	.text;																				\
	.global TEST_FUNC_NAME;												\
	.global TEST_FUNC_RET;												\
TEST_FUNC_NAME:																	\
 la a0, .test_name;															\
 li	a2, 0x80001000;															\
.prname_next:																		\
 lb	a1,0(a0);																		\
 addi	a0,a0,1;																	\
 beq	a1,zero,.prname_done;											\
 sb	a1,0(a2);																		\
 jal	zero,.prname_next;												\
.test_name:																			\
 .ascii TEST_FUNC_TXT;													\
 .byte 0x00;																		\
 .balign 4;																			\
.prname_done:																		\
 addi	a1,zero,'.';															\
 sb	a1,0(a2);																		\
 sb	a1,0(a2);

#define RVTEST_PASS															\
	li	  a0,0x80001000; 													\
	addi	a1,zero,'O';														\
	addi	a2,zero,'K';														\
	addi	a3,zero,'\n';														\
	sb	a1,0(a0);																	\
	sb	a2,0(a0);																	\
	sb	a3,0(a0);																	\
	jal	zero, TEST_FUNC_RET;

#define RVTEST_FAIL															\
	li	a0,0x80001000; 														\
	addi	a1,zero,'E';														\
	addi	a2,zero,'R';														\
	addi	a3,zero,'O';														\
	addi	a4,zero,'\n';														\
	sb	a1,0(a0);																	\
	sb	a2,0(a0);																	\
	sb	a2,0(a0);																	\
	sb	a3,0(a0);																	\
	sb	a2,0(a0);																	\
	sb	a4,0(a0);																	\
	jal zero, TEST_FUNC_RET;

#define RVTEST_CODE_END
#define RVTEST_DATA_BEGIN .balign 4;
#define RVTEST_DATA_END

#endif
