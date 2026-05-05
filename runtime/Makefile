# Makefile for Fuji Runtime Library
# Use a POSIX shell so `mkdir -p` / `rm -rf` work when GNU Make runs on Windows (cmd.exe does not).
SHELL := sh.exe
.SHELLFLAGS := -ec

CC = gcc
CFLAGS = -O3 -fPIC -std=c11
AR = ar
ARFLAGS = rcs

SRCDIR = src
OBJDIR = obj
# Output: runtime/libfuji_runtime.a (linked by nativebuild)
LIBDIR = .

SRCS = $(SRCDIR)/value.c $(SRCDIR)/object.c $(SRCDIR)/gc.c $(SRCDIR)/fuji_runtime.c
OBJS = $(patsubst $(SRCDIR)/%.c,$(OBJDIR)/%.o,$(SRCS))
TARGET = $(LIBDIR)/libfuji_runtime.a

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(OBJS)
	@mkdir -p $(LIBDIR)
	$(AR) $(ARFLAGS) $@ $^

$(OBJDIR)/%.o: $(SRCDIR)/%.c
	@mkdir -p $(OBJDIR)
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -rf $(OBJDIR) $(TARGET)
