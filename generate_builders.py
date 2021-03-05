import datetime
from buildbot.plugins import *

build_lock = util.WorkerLock("worker_builds",
                             maxCount=9999)

targets = [
#            ( "acorn32", "" ),
#            ( "algor", "" ),
#            ( "alpha", "" ),
#            ( "amd64", "" ),
#            ( "amiga", "" ),
#            ( "amigappc", "" ),
#            ( "arc", "" ),
#            ( "atari", "" ),
#            ( "bebox", "" ),
#            ( "cats", "" ),
#            ( "cesfic", "" ),
#            ( "cobalt", "" ),
#            ( "dreamcast", "" ),
#            ( "emips", "" ),
#            ( "epoc32", "" ),
#            ( "evbarm", "aarch64" ),
#            ( "evbarm", "aarch64eb" ),
#            ( "evbarm", "earm" ),
#            ( "evbarm", "earmeb" ),
#            ( "evbarm", "earmhf" ),
#            ( "evbarm", "earmhfeb" ),
#            ( "evbarm", "earmv5" ),
#            ( "evbarm", "earmv5eb" ),
#            ( "evbarm", "earmv5hf" ),
#            ( "evbarm", "earmv5hfeb" ),
#            ( "evbarm", "earmv6hf" ),
#            ( "evbarm", "earmv6hfeb" ),
#            ( "evbarm", "earmv7hf" ),
#            ( "evbarm", "earmv7hfeb" ),
#            ( "evbmips", "mips64eb" ),
#            ( "evbmips", "mips64el" ),
#            ( "evbmips", "mipseb" ),
#            ( "evbmips", "mipsel" ),
#            ( "evbppc", "" ),
#            ( "evbppc", "powerpc64" ),
#            ( "evbsh3", "sh3eb" ),
#            ( "evbsh3", "sh3el" ),
#            ( "ews4800mips", "" ),
#            ( "hp300", "" ),
#            ( "hpcarm", "" ),
#            ( "hpcmips", "" ),
#            ( "hpcsh", "" ),
#            ( "hppa", "" ),
            ( "i386", "" ),
#            ( "ia64", "" ),
#            ( "ibmnws", "" ),
#            ( "iyonix", "" ),
#            ( "landisk", "" ),
#            ( "luna68k", "" ),
#            ( "mac68k", "" ),
#            ( "macppc", "" ),
#            ( "mipsco", "" ),
#            ( "mmeye", "" ),
#            ( "mvme68k", "" ),
#            ( "mvmeppc", "" ),
#            ( "netwinder", "" ),
#            ( "news68k", "" ),
#            ( "newsmips", "" ),
#            ( "next68k", "" ),
#            ( "ofppc", "" ),
#            ( "pmax", "" ),
#            ( "prep", "" ),
#            ( "rs6000", "" ),
#            ( "sandpoint", "" ),
#            ( "sgimips", "" ),
#            ( "shark", "" ),
#            ( "sparc64", "" ),
#            ( "sparc", "" ),
#            ( "sun2", "" ),
#            ( "sun3", "" ),
#            ( "vax", "" ),
#            ( "x68k", "" ),
#            ( "zaurus", "" ),
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

def is_test_target(target):
    arch = target[TARGET_ARCH]
    machine = target[TARGET_MACHINE]
    if (arch == "amd64" or
        arch == "i386" or
        arch == "sparc64" or
        arch == "sparc" or
        arch == "pmax" or
        arch == "hpcmips" or
        machine == "aarch64" or
        machine == "earmv7hf"):
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

def to_builder(targets, buildtype, branchname):
    if buildtype != "":
        build_name = branchname + "-" + buildtype
        tags = [buildtype, branchname]
    else:
        build_name = branchname
        tags = [branchname]

    def target_name(target):
        arch = target[TARGET_ARCH]
        machine = target[TARGET_MACHINE]
        if machine != "":
            return arch + "-" + machine
        else:
            return arch

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
            if buildtype == "RELEASE":
                return ["-V", "NETBSD_OFFICIAL_RELEASE=yes"]
            return []

        return " ".join(["rm", "-rf", "../build/*", ";", 
                 "../src/build.sh",
                 "-O", "$PWD",
                 "-j", "6", # XXX ncpu "$(/sbin/sysctl -n hw.ncpuonline)",
                 "-B", "$(date -r $(cd ../src; git show -s --format=%ct) +%Y%m%d%H%MZ)",
                 "-R", "$HOME/releasedir/" + build_name  + "/$(date -r $(cd ../src; git show -s --format=%ct) +%Y%m%d%H%MZ)/",
                 "-U", "-P", "-N0", "-V", "TMPDIR=/tmp",
                 "-V", "MKDEBUG=yes", "-V", "BUILD=yes"]
                + x_flags(target) + target_flags(target) +
                buildtype_flags(buildtype) + build_target(target, branchname) +
                [";", "build_status=$?"] +
                [";", "rm", "-rf", "../build/*"] +
                [";", "exit", "$build_status"])

    factory = util.BuildFactory()
    factory.addStep(steps.Git(
                haltOnFailure=True,
                logEnviron=False,
                repourl='https://github.com/netbsd/src.git',
                branch=branchname,
                mode='incremental',
                codebase='src',
                workdir="src",
                timeout=12000
            ))
    factory.addStep(steps.Git(
                haltOnFailure=True,
                logEnviron=False,
                repourl='https://github.com/netbsd/xsrc.git',
                branch=branchname,
                mode='incremental',
                codebase='xsrc',
                workdir="xsrc",
                timeout=12000
            ))
    factory.addStep(steps.ShellCommand(
                haltOnFailure=True,
                logEnviron=False,
                name="clean releasedir before",
                description="cleaning release directory - before",
                descriptionDone="clean releasedir",
                command="rm -rf $HOME/releasedir"
            ))

    for target in targets:
        factory.addStep(steps.ShellCommand(
                    haltOnFailure=False,
                    logEnviron=False,
                    name="build " + build_name + " " + target_name(target),
                    description="building src",
                    descriptionDone="build done",
                    command=build_command(target, buildtype),
                    locks=[build_lock.access('exclusive')],
                    workdir="build"
                ))

    factory.addStep(steps.ShellCommand(
                haltOnFailure=True,
                logEnviron=False,
                name="uploading releasedir",
                description="uploading release directory",
                descriptionDone="upload releasedir",
                command="rsync -avr $HOME/releasedir $HOME/releasedir-target-upload/"
            ))

    def buildtype_is_tested(buildtype):
        if (buildtype == "lint" or
            buildtype == "LLVM"):
            return False
        return True

    if buildtype_is_tested(buildtype):
        test_targets = [target for target in targets if is_test_target(target)]
        for target in test_targets:
            factory.addStep(steps.ShellCommand(
                        haltOnFailure=False,
                        logEnviron=False,
                        name="testing " + build_name + " " + target_name(target),
                        description="testing " + build_name + " " + target_name(target),
                        descriptionDone="testing done",
                        command="anita test --vmm-args \"-accel nvmm\" --memory-size 512M --workdir workdir-tests/ $HOME/releasedir/" + build_name + "/$(date -r $(cd ../src; git show -s --format=%ct) +%Y%m%d%H%MZ)/" + target_name(target) + "/",
                        timeout=6000,
                    ))
            factory.addStep(steps.ShellCommand(
                        haltOnFailure=True,
                        logEnviron=False,
                        name="test output (XSL) " + build_name + " " + target_name(target),
                        description="test output (XSL) " + build_name + " " + target_name(target),
                        descriptionDone="output done",
                        command="cat workdir-tests/atf/tests-results.xsl",
                    ))
            factory.addStep(steps.ShellCommand(
                        haltOnFailure=True,
                        logEnviron=False,
                        name="test output (XML) " + build_name + " " + target_name(target),
                        description="test output (XML) " + build_name + " " + target_name(target),
                        descriptionDone="output done",
                        command="cat workdir-tests/atf/test.xml",
                    ))
            factory.addStep(steps.ShellCommand(
                        haltOnFailure=True,
                        logEnviron=False,
                        name="delete test workdir " + build_name + " " + target_name(target),
                        description="deleting test workdir " + build_name + " " + target_name(target),
                        descriptionDone="deleting test workdir done",
                        command="rm -rf workdir-tests",
                    ))

    return util.BuilderConfig(name=build_name,
                              workernames=["worker1"],
                              factory=factory,
                              tags=tags,
                              ), build_name

def generate_stable_builders():
    netbsd_8_builder = to_builder(
            [target for target in targets if is_8_target(target)],
            buildtype="",
            branchname="netbsd-8")
    netbsd_9_builder = to_builder(
            [target for target in targets if is_9_target(target)],
            buildtype="",
            branchname="netbsd-9")

    return [netbsd_8_builder, netbsd_9_builder]


def generate_head_builders():
    llvm_builder = to_builder(
            [target for target in targets if is_llvm_target(target)],
            buildtype="LLVM",
            branchname="HEAD")
    lint_builder = to_builder(
            [target for target in targets if is_lint_target(target)],
            buildtype="lint",
            branchname="HEAD")
    HEAD_builder = to_builder(
            [target for target in targets if is_HEAD_target(target)],
            buildtype="",
            branchname="HEAD")

    return [HEAD_builder, llvm_builder, lint_builder]

def generate_release_builders():
    netbsd_8_RELEASE_builder = to_builder(
            [target for target in targets if is_8_target(target)],
            buildtype="RELEASE",
            branchname="netbsd-8")
    netbsd_9_RELEASE_builder = to_builder(
            [target for target in targets if is_9_target(target)],
            buildtype="RELEASE",
            branchname="netbsd-9")

    return [netbsd_8_RELEASE_builder, netbsd_9_RELEASE_builder]
