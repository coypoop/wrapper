import datetime
from buildbot.plugins import *

def generate_builders():
    builders = []
    targets = [
                ( "acorn32", "" ),
                ( "algor", "" ),
                ( "alpha", "" ),
                ( "amd64", "" ),
                ( "amiga", "" ),
                ( "amigappc", "" ),
                ( "arc", "" ),
                ( "atari", "" ),
                ( "bebox", "" ),
                ( "cats", "" ),
                ( "cesfic", "" ),
                ( "cobalt", "" ),
                ( "dreamcast", "" ),
                ( "emips", "" ),
                ( "epoc32", "" ),
                ( "evbarm", "aarch64" ),
                ( "evbarm", "aarch64eb" ),
                ( "evbarm", "earm" ),
                ( "evbarm", "earmeb" ),
                ( "evbarm", "earmhf" ),
                ( "evbarm", "earmhfeb" ),
                ( "evbarm", "earmv5" ),
                ( "evbarm", "earmv5eb" ),
                ( "evbarm", "earmv5hf" ),
                ( "evbarm", "earmv5hfeb" ),
                ( "evbarm", "earmv6hf" ),
                ( "evbarm", "earmv6hfeb" ),
                ( "evbarm", "earmv7hf" ),
                ( "evbarm", "earmv7hfeb" ),
                ( "evbmips", "mips64eb" ),
                ( "evbmips", "mips64el" ),
                ( "evbmips", "mipseb" ),
                ( "evbmips", "mipsel" ),
                ( "evbppc", "" ),
                ( "evbppc", "powerpc64" ),
                ( "evbsh3", "sh3eb" ),
                ( "evbsh3", "sh3el" ),
                ( "ews4800mips", "" ),
                ( "hp300", "" ),
                ( "hpcarm", "" ),
                ( "hpcmips", "" ),
                ( "hpcsh", "" ),
                ( "hppa", "" ),
                ( "i386", "" ),
                ( "ia64", "" ),
                ( "ibmnws", "" ),
                ( "iyonix", "" ),
                ( "landisk", "" ),
                ( "luna68k", "" ),
                ( "mac68k", "" ),
                ( "macppc", "" ),
                ( "mipsco", "" ),
                ( "mmeye", "" ),
                ( "mvme68k", "" ),
                ( "mvmeppc", "" ),
                ( "netwinder", "" ),
                ( "news68k", "" ),
                ( "newsmips", "" ),
                ( "next68k", "" ),
                ( "ofppc", "" ),
                ( "pmax", "" ),
                ( "prep", "" ),
                ( "rs6000", "" ),
                ( "sandpoint", "" ),
                ( "sgimips", "" ),
                ( "shark", "" ),
                ( "sparc64", "" ),
                ( "sparc", "" ),
                ( "sun2", "" ),
                ( "sun3", "" ),
                ( "vax", "" ),
                ( "x68k", "" ),
                ( "zaurus", "" ),
    ]
    TARGET_ARCH = 0
    TARGET_MACHINE = 1

    def is_lint_target(target):
        arch = target[TARGET_ARCH]
        if (arch == "amd64" or
            arch == "i386" or
            arch == "sparc64" or
            arch == "sparc"):
           return True
        return False

    def is_llvm_target(target):
        arch = target[TARGET_ARCH]
        machine = target[TARGET_MACHINE]
        if (arch == "amd64" or
            arch == "i386" or
            arch == "sparc64" or
            arch == "sparc" or
            arch == "macppc" or
            arch == "evbppc" or
            machine == "aarch64" or
            machine == "earmv6hf"):
            return True
        return False

    def target_in_branch(target, branch):
        arch = target[TARGET_ARCH]
        machine = target[TARGET_MACHINE]
        if (machine == "aarch64"):
            return (branch >= 9)
        if (machine == "aarch64eb"):
            return (branch >= 10)
        if (machine == "earm" or
            machine == "earmeb"):
            return (branch < 10)
        if (machine == "earmhf" or
            machine == "earmhfeb"):
            return (branch == 9)
        if arch == "ia64":
            return (branch >= 9)
        if machine == "powerpc64":
            return (branch >= 9)
        if (machine == "earmv5" or
            machine == "earmv5eb" or
            machine == "earmv5" or
            machine == "earmv5hf" or
            machine == "earmv5hfeb" or
            machine == "earmv5hfeb"):
            return (branch >= 10)
        if (machine == "earmv6hfeb"):
            return (branch >= 10)
        return True

    def is_8_target(target):
        return target_in_branch(target, 8)
    def is_9_target(target):
        return target_in_branch(target, 9)
    def is_HEAD_target(target):
        return target_in_branch(target, 10)

    def to_builder(target, buildtype, branchname):
        def build_target(target, branchname):
            arch = target[TARGET_ARCH]
            machine = target[TARGET_MACHINE]
            if (arch == "evbarm" and
                branchname == "HEAD"):
                if (machine == "aarch64" or
                    machine == "aarch64eb"):
                    return ["release", "install-image", "iso-image"]
                if (machine == "earmv7hf" or
                    machine == "earmv7hfeb"):
                    return ["release", "install-image"]
                return ["release"]
            if (arch == "mac68k"):
                return ["release"]
            if (arch == "i386"):
                return ["release", "iso-image", "install-image"]
            if (arch == "amd64"):
                if (branchname == "HEAD"):
                    return ["release", "iso-image", "install-image", "live-image"]
                return ["release", "iso-image", "install-image"]
            return ["release", "iso-image"]

        def build_command(target, buildtype):
            def x_flags(target):
                arch = target[TARGET_ARCH]
                if (arch == "evbppc" or
                    arch == "rs6000" or
                    arch == "sun2"):
                    return []
                return ["-x"]

            def target_flags(target):
                arch = target[TARGET_ARCH]
                machine = target[TARGET_MACHINE]

                target_string = ["-m", arch]
                if machine != "":
                    target_string = target_string + ["-a", machine]
                return target_string

            def buildtype_flags(buildtype):
                if buildtype == "LLVM":
                    return ["-V", "MKLLVM=yes", "-V", "HAVE_LLVM=yes", "-V", "MKGCC=no"]
                if buildtype == "lint":
                    return ["-V", "MKLINT=yes"]
                return []

            return (["./wrapper.sh",
                     "-U", "-P", "-N0", "-V", "TMPDIR=/tmp",
                     "-V", "MKDEBUG=yes", "-V", "BUILD=yes"]
                    + x_flags(target) + target_flags(target) + buildtype_flags(buildtype) + build_target(target, branchname))

        factory = util.BuildFactory()
        factory.addStep(steps.Git(
                    haltOnFailure=True,
                    logEnviron=False,
                    repourl='https://github.com/netbsd/src.git',
                    branch=branchname,
                    mode='incremental',
                    codebase='src',
                    retry=(5, 3),
                    workdir="src"
                ))
        factory.addStep(steps.Git(
                    haltOnFailure=True,
                    logEnviron=False,
                    repourl='https://github.com/netbsd/xsrc.git',
                    branch=branchname,
                    mode='incremental',
                    codebase='xsrc',
                    retry=(5, 3),
                    workdir="xsrc"
                ))
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=True,
                    logEnviron=False,
                    name="clean-obj-before",
                    description="cleaning obj directory - before",
                    descriptionDone="clean obj",
                    command=["rm", "-rf", "../build"]
                ))
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=True,
                    logEnviron=False,
                    name="fetch-wrapper",
                    description="fetch wrapper",
                    descriptionDone="fetched wrapper",
                    command=["ftp", "https://raw.githubusercontent.com/coypoop/wrapper/main/wrapper.sh"],
                    workdir="build"
                ))
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=True,
                    logEnviron=False,
                    name="chmod-wrapper",
                    description="chmod wrapper",
                    descriptionDone="chmodded wrapper",
                    command=["chmod", "+x", "wrapper.sh"],
                    workdir="build"
                ))
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=False,
                    logEnviron=False,
                    name="build",
                    description="building src",
                    descriptionDone="build done",
                    command=build_command(target, buildtype),
                    workdir="build"
                ))
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=True,
                    logEnviron=False,
                    name="clean-obj-after",
                    description="cleaning obj directory - after",
                    descriptionDone="clean obj",
                    command=["rm", "-rf", "../build"]
                ))
        arch = target[TARGET_ARCH]
        machine = target[TARGET_MACHINE]
        build_name = arch + "-" + machine + "-" + buildtype + "-" + branchname
        tags = [arch]
        if buildtype != "":
            tags.append(buildtype)
        if machine != "":
            tags.append(machine)

        return util.BuilderConfig(name=build_name,
                                  workernames=["worker1"],
                                  factory=factory,
                                  tags=tags)

    def to_builder_8(target):
        return to_builder(target, buildtype="", branchname="netbsd-8")
    def to_builder_9(target):
        return to_builder(target, buildtype="", branchname="netbsd-9")
    def to_builder_HEAD(target):
        return to_builder(target, buildtype="", branchname="trunk")
    def to_builder_HEAD_llvm(target):
        return to_builder(target, buildtype="LLVM", branchname="trunk")
    def to_builder_HEAD_lint(target):
        return to_builder(target, buildtype="lint", branchname="trunk")

    llvm_targets = filter(is_llvm_target, targets)
    lint_targets = filter(is_lint_target, targets)
    netbsd_8_targets = filter(is_8_target, targets)
    netbsd_9_targets = filter(is_9_target, targets)
    netbsd_HEAD_targets = filter(is_HEAD_target, targets)

    llvm_builders = list(map(to_builder_HEAD_llvm, llvm_targets))
    lint_builders = list(map(to_builder_HEAD_lint, lint_targets))
    netbsd_8_builders = list(map(to_builder_8, netbsd_8_targets))
    netbsd_9_builders = list(map(to_builder_9, netbsd_9_targets))
    netbsd_HEAD_builders = list(map(to_builder_HEAD, netbsd_HEAD_targets))

    return llvm_builders + lint_builders + netbsd_8_builders + netbsd_9_builders + netbsd_HEAD_builders
